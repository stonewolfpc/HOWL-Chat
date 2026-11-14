package contentimport

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

const (
	maxSourceDocumentBytes = 8 << 20
	maxSourceLoreBytes     = 4 << 20
	maxSourceMediaBytes    = 8 << 20
)

var (
	worldLoreNames    = []string{"world_lore.json", "automated_world_lore.json"}
	scenarioLoreNames = []string{"scenario_lore.json", "automated_scenario_lore.json"}
)

func (s *Service) documentName() string {
	if s.category == library.Scenarios {
		return "scenario.json"
	}
	return "world.json"
}

func (s *Service) loreNames() []string {
	if s.category == library.Scenarios {
		return append([]string(nil), scenarioLoreNames...)
	}
	return append([]string(nil), worldLoreNames...)
}

func (s *Service) loadBrowserSources(records []map[string]interface{}) ([]sourceInfo, error) {
	out := make([]sourceInfo, 0, len(records))
	for i, record := range records {
		if record == nil {
			out = append(out, sourceInfo{
				key:          fmt.Sprintf("browser:%d", i),
				kind:         "browser",
				blockReason:  "invalid_document_id",
				blockMessage: "The supplied browser record is not an object.",
			})
			continue
		}
		normalized, err := normalizeMap(record)
		if err != nil {
			return nil, fmt.Errorf("browser record %d: %w", i, err)
		}
		rawHash := hashJSON(normalized)
		src := sourceInfo{
			key:          fmt.Sprintf("browser:%d", i),
			kind:         "browser",
			raw:          normalized,
			body:         normalized,
			rawHash:      rawHash,
			docHash:      hashJSON(normalized),
			lore:         map[string][]byte{},
			sourceHashes: map[string]string{"record.json": rawHash},
		}
		src.id = strings.TrimSpace(stringField(normalized, "id"))
		src.name = displayName(normalized)
		if inlineMediaIn(normalized) {
			src.blockReason = "inline_media_unmigrated"
			src.blockMessage = "The browser record still contains inline media bytes; image-byte migration is not part of this import."
		}
		out = append(out, src)
	}
	return out, nil
}

func (s *Service) loadLegacySources(roots []string, destinationRoot string) ([]sourceInfo, []string, error) {
	var out []sourceInfo
	var errorsOut []string
	for rootIndex, rawRoot := range roots {
		root := strings.TrimSpace(rawRoot)
		if root == "" {
			return nil, nil, fmt.Errorf("legacy source root is empty")
		}
		abs, err := filepath.Abs(root)
		if err != nil {
			return nil, nil, fmt.Errorf("legacy source root %q: %w", root, err)
		}
		info, err := os.Lstat(abs)
		if err != nil {
			return nil, nil, fmt.Errorf("legacy source root %q: %w", abs, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return nil, nil, fmt.Errorf("legacy source root %q is not a regular directory", abs)
		}
		abs = canonicalPath(abs)
		if destinationRoot != "" {
			destination := canonicalPath(destinationRoot)
			if library.IsUnderDir(abs, destination) || library.IsUnderDir(destination, abs) {
				out = append(out, sourceInfo{
					key:          fmt.Sprintf("legacy:%d", rootIndex),
					kind:         "legacy",
					root:         abs,
					blockReason:  "source_destination_overlap",
					blockMessage: "The legacy source root overlaps the destination root; no write is planned.",
				})
				continue
			}
		}
		if view, registered := s.registeredRootOverlap(abs); registered {
			out = append(out, sourceInfo{
				key:          fmt.Sprintf("legacy:%d", rootIndex),
				kind:         "legacy",
				root:         abs,
				blockReason:  "source_registered",
				blockMessage: "The selected legacy source is a registered content root or overlaps one: " + view.Path,
			})
			continue
		}
		entries, err := os.ReadDir(abs)
		if err != nil {
			return nil, nil, fmt.Errorf("read legacy source root %q: %w", abs, err)
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
		for _, entry := range entries {
			folder := entry.Name()
			full := filepath.Join(abs, folder)
			info, err := os.Lstat(full)
			if err != nil {
				errorsOut = append(errorsOut, fmt.Sprintf("%s: %v", full, err))
				continue
			}
			if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
				errorsOut = append(errorsOut, fmt.Sprintf("%s: unsupported source entry (not a regular directory)", full))
				continue
			}
			src, err := s.readLegacyFolder(abs, folder, full, fmt.Sprintf("legacy:%d:%s", rootIndex, folder))
			if err != nil {
				errorsOut = append(errorsOut, fmt.Sprintf("%s: %v", full, err))
				continue
			}
			out = append(out, src)
		}
	}
	return out, errorsOut, nil
}

func (s *Service) readLegacyFolder(root, folder, full, key string) (sourceInfo, error) {
	src := sourceInfo{
		key:          key,
		kind:         "legacy",
		root:         root,
		folder:       folder,
		name:         folder,
		lore:         map[string][]byte{},
		sourceHashes: map[string]string{},
	}
	docName := s.documentName()
	allowed := map[string]bool{docName: true, "media.json": true}
	for _, name := range s.loreNames() {
		allowed[name] = true
	}
	entries, err := os.ReadDir(full)
	if err != nil {
		return src, fmt.Errorf("read legacy folder: %w", err)
	}
	var documentPath string
	var documentFiles []string
	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(full, name)
		info, err := os.Lstat(path)
		if err != nil {
			return src, err
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			src.unsupported = append(src.unsupported, name)
			continue
		}
		if !allowed[name] {
			src.unsupported = append(src.unsupported, name)
			continue
		}
		if name == docName {
			documentFiles = append(documentFiles, name)
		}
		_ = path
	}
	sort.Strings(documentFiles)
	if len(documentFiles) != 1 {
		src.blockReason = "unsupported_source_file"
		src.blockMessage = fmt.Sprintf("Expected exactly one %s in the folder; found %d.", docName, len(documentFiles))
		return src, nil
	}
	documentPath = filepath.Join(full, documentFiles[0])
	payload, err := readLimited(documentPath, maxSourceDocumentBytes)
	if err != nil {
		return src, fmt.Errorf("read document: %w", err)
	}
	src.sourceHashes[documentFiles[0]] = hashBytes(payload)
	raw, err := decodeObject(payload)
	if err != nil {
		return src, fmt.Errorf("decode document: %w", err)
	}
	src.raw = raw
	body := raw
	outerID := strings.TrimSpace(stringField(raw, "id"))
	if inner, ok := raw["document"].(map[string]interface{}); ok && strings.TrimSpace(stringField(raw, "format")) != "" {
		body = inner
	}
	if len(body) == 0 {
		src.blockReason = "invalid_document_id"
		src.blockMessage = "The legacy document body is empty."
		return src, nil
	}
	src.body = body
	src.id = strings.TrimSpace(stringField(body, "id"))
	if src.id == "" {
		src.id = outerID
	}
	src.name = displayName(body)
	if src.name == folder {
		src.name = displayName(raw)
	}
	src.rawHash = hashJSON(body)
	src.docHash = src.rawHash
	for _, name := range s.loreNames() {
		path := filepath.Join(full, name)
		if _, err := os.Lstat(path); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return src, err
		}
		data, err := readLimited(path, maxSourceLoreBytes)
		if err != nil {
			return src, fmt.Errorf("read lore %q: %w", name, err)
		}
		src.lore[name] = data
		src.sourceHashes[name] = hashBytes(data)
	}
	mediaPath := filepath.Join(full, "media.json")
	if _, err := os.Lstat(mediaPath); err == nil {
		data, err := readLimited(mediaPath, maxSourceMediaBytes)
		if err != nil {
			return src, fmt.Errorf("read media.json: %w", err)
		}
		normalized, err := worldstore.NormalizeMediaManifest(data)
		if err != nil {
			switch {
			case errors.Is(err, worldstore.ErrMediaInline):
				src.blockReason = "inline_media_unmigrated"
				src.blockMessage = "media.json contains image bytes; only reference-only media is preserved."
			case errors.Is(err, worldstore.ErrMediaUnsupported):
				src.blockReason = "unsupported_media_manifest"
				src.blockMessage = "media.json is a structurally valid but unsupported reference-only manifest: " + err.Error()
			default:
				src.blockReason = "malformed_media_manifest"
				src.blockMessage = "media.json is not a valid reference-only manifest: " + err.Error()
			}
			return src, nil
		}
		src.media = normalized
		src.sourceHashes["media.json"] = hashBytes(normalized)
	}
	if inlineMediaIn(body) {
		src.blockReason = "inline_media_unmigrated"
		src.blockMessage = "The legacy document contains inline media bytes."
	}
	if len(src.unsupported) > 0 {
		src.unsupported = append([]string(nil), src.unsupported...)
		sort.Strings(src.unsupported)
		src.blockReason = "unsupported_source_file"
		src.blockMessage = "Unsupported source files: " + strings.Join(src.unsupported, ", ")
	}
	return src, nil
}

func (s *Service) prepareItem(src sourceInfo, destinationID string, evidence string) (preparedItem, string, string) {
	document, err := normalizeMap(src.body)
	if err != nil {
		return preparedItem{}, "invalid_document_id", err.Error()
	}
	document["id"] = destinationID
	references := s.normalizeReferences(document)
	origin := worldstore.Origin{Kind: "imported"}
	if strings.TrimSpace(src.id) != "" && src.id != destinationID {
		origin.LegacyIDs = []string{src.id}
	}
	if strings.TrimSpace(src.name) != "" && src.name != destinationID {
		origin.LegacyName = src.name
	}
	lore := map[string][]byte{}
	for _, filename := range s.loreNames() {
		raw, ok := src.lore[filename]
		if !ok {
			continue
		}
		automated := strings.HasPrefix(filename, "automated_")
		data, reason, message := prepareLoreFile(raw, destinationID, automated)
		if reason != "" {
			return preparedItem{}, reason, message
		}
		lore[filename] = data
	}
	createdAt := strings.TrimSpace(stringField(document, "createdAt"))
	if createdAt == "" {
		createdAt = strings.TrimSpace(stringField(document, "created_at"))
	}
	return preparedItem{
		source:     src,
		document:   document,
		origin:     origin,
		createdAt:  createdAt,
		lore:       lore,
		media:      append([]byte(nil), src.media...),
		references: references,
	}, "", evidence
}

func prepareLoreFile(raw []byte, destinationID string, automated bool) ([]byte, string, string) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var doc map[string]interface{}
	if err := dec.Decode(&doc); err != nil {
		return nil, "unsupported_source_file", "Lore file is not a JSON object."
	}
	entries, ok := doc["entries"].([]interface{})
	if !ok && doc["entries"] != nil {
		return nil, "unsupported_source_file", "Lore entries are not an array."
	}
	for i, value := range entries {
		entry, ok := value.(map[string]interface{})
		if !ok {
			return nil, "unsupported_source_file", fmt.Sprintf("Lore entry %d is not an object.", i)
		}
		for _, key := range []string{"ownerId", "owner_id", "owner"} {
			if owner, ok := entry[key].(string); ok && strings.TrimSpace(owner) != "" && strings.TrimSpace(owner) != destinationID {
				return nil, "unmapped_owner", fmt.Sprintf("Lore entry %d has a different owner; it was not rewritten.", i)
			}
		}
		ownerID, hasOwnerID := entry["ownerId"].(string)
		if !hasOwnerID || strings.TrimSpace(ownerID) == "" {
			entry["ownerId"] = destinationID
		}
		if automated {
			entry["automated"] = true
			entry["originKind"] = "legacy_generated_unattributed"
		}
	}
	doc["entries"] = entries
	payload, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, "unsupported_source_file", "Lore file could not be re-encoded."
	}
	if len(payload) > maxSourceLoreBytes {
		return nil, "unsupported_source_file", "Lore file exceeds the import size limit."
	}
	return payload, "", ""
}

func semanticLoreHash(raw []byte) string {
	if len(bytes.TrimSpace(raw)) == 0 {
		return hashJSON(map[string]interface{}{"entries": []interface{}{}})
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var doc map[string]interface{}
	if err := dec.Decode(&doc); err != nil {
		return ""
	}
	entries, _ := doc["entries"].([]interface{})
	if entries == nil {
		entries = []interface{}{}
	}
	return hashJSON(map[string]interface{}{"entries": entries})
}

func inlineMediaIn(value interface{}) bool {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, item := range typed {
			switch key {
			case "data", "base64", "blob", "dataUrl", "data_url", "imageData", "image_data":
				if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
					return true
				}
			case "image", "images", "thumbnail", "avatar", "media", "legacyRef":
				if inlineMediaIn(item) {
					return true
				}
			default:
				if inlineMediaIn(item) {
					return true
				}
			}
		}
	case []interface{}:
		for _, item := range typed {
			if inlineMediaIn(item) {
				return true
			}
		}
	case string:
		trimmed := strings.TrimSpace(typed)
		lower := strings.ToLower(trimmed)
		return strings.HasPrefix(lower, "data:image/") || (strings.Contains(lower, ";base64,") && len(trimmed) > 64)
	}
	return false
}

func canonicalPath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	return filepath.Clean(path)
}

func (s *Service) registeredRootOverlap(path string) (library.View, bool) {
	if s.reg == nil {
		return library.View{}, false
	}
	for _, view := range s.reg.List("") {
		registered := canonicalPath(view.Path)
		if library.IsUnderDir(path, registered) || library.IsUnderDir(registered, path) {
			return view, true
		}
	}
	return library.View{}, false
}

func displayName(doc map[string]interface{}) string {
	for _, key := range []string{"title", "name"} {
		if value := strings.TrimSpace(stringField(doc, key)); value != "" {
			return value
		}
	}
	return ""
}

func stringField(doc map[string]interface{}, key string) string {
	value, _ := doc[key].(string)
	return value
}

func normalizeMap(raw map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return decodeObject(data)
}

func decodeObject(data []byte) (map[string]interface{}, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var out map[string]interface{}
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("JSON value is not an object")
	}
	return out, nil
}

func readLimited(path string, max int64) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Size() > max {
		return nil, fmt.Errorf("file exceeds %d bytes", max)
	}
	return os.ReadFile(path)
}

func hashJSON(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return hashBytes(data)
}

func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
