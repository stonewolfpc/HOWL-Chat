package worldstore

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/library"
)

const (
	// ImportStagingDir is a reserved, non-item directory under a destination root.
	// Its first segment is not a valid item ID, so Lookup cannot address it.
	ImportStagingDir = ".howl-import-staging"
	importOwnerFile  = "owner.json"
	importOwnerMagic = "howl-import-staging-v1"
)

var (
	// ErrImportStagingExists reports a reused operation/item staging path.
	ErrImportStagingExists = errors.New("import staging path already exists")
	// ErrImportConflict reports staging/owner evidence that does not match.
	ErrImportConflict = errors.New("import staging ownership does not match")
)

// ImportItem is one canonical item prepared for staged import. The import
// orchestrator supplies bytes; this package owns all writes into the root.
type ImportItem struct {
	ID                   string
	DestinationLibraryID string
	Document             map[string]interface{}
	Origin               Origin
	CreatedAt            string
	Lore                 map[string][]byte
	Media                []byte
	PlanHash             string
}

// LoadStaging re-reads an owned staging item for explicit resume. It accepts
// only a complete owner marker matching the operation, item, library and plan.
func (s *Store) LoadStaging(operationID, itemID, libraryID, planHash string) (StagedImport, error) {
	if s == nil || s.reg == nil {
		return StagedImport{}, fmt.Errorf("%w: content library registry is not initialised", ErrUnavailable)
	}
	if !validID.MatchString(operationID) || !validID.MatchString(itemID) {
		return StagedImport{}, fmt.Errorf("%w: invalid staging identity", ErrInvalidDocument)
	}
	view, ok := s.registrationView(libraryID)
	if !ok || view.Category != s.category {
		return StagedImport{}, ErrNoDestination
	}
	staged := StagedImport{
		OperationID: operationID,
		ItemID:      itemID,
		LibraryID:   view.ID,
		Root:        view.Path,
		StagingPath: filepath.Join(view.Path, ImportStagingDir, operationID, itemID),
		OwnerPath:   filepath.Join(view.Path, ImportStagingDir, operationID, itemID+"."+importOwnerFile),
		FinalPath:   filepath.Join(view.Path, itemID),
		PlanHash:    planHash,
	}
	data, err := os.ReadFile(staged.OwnerPath)
	if err != nil {
		return StagedImport{}, fmt.Errorf("%w: read staging ownership: %v", ErrImportConflict, err)
	}
	var owner importOwner
	if err := json.Unmarshal(data, &owner); err != nil {
		return StagedImport{}, fmt.Errorf("%w: decode staging ownership: %v", ErrImportConflict, err)
	}
	staged.DocumentHash = owner.DocumentHash
	staged.LoreHashes = owner.LoreHashes
	staged.MediaHash = owner.MediaHash
	staged.Revision = owner.Revision
	if err := s.validateStagedPaths(staged); err != nil {
		return StagedImport{}, err
	}
	if err := s.verifyImportOwner(staged); err != nil {
		return StagedImport{}, err
	}
	if err := s.verifyStagedFiles(staged); err != nil {
		return StagedImport{}, fmt.Errorf("%w: staged files changed: %v", ErrImportConflict, err)
	}
	return staged, nil
}

// VerifyImported proves an already-committed item still matches its journaled
// revision and every journaled file hash. It changes nothing.
func (s *Store) VerifyImported(id, expectedRevision string, expectedHashes map[string]string) (bool, error) {
	if s == nil || s.reg == nil {
		return false, fmt.Errorf("%w: content library registry is not initialised", ErrUnavailable)
	}
	if !validID.MatchString(id) {
		return false, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	if len(expectedHashes) == 0 {
		return false, fmt.Errorf("%w: no journaled hashes", ErrInvalidDocument)
	}
	writeMu.Lock()
	defer writeMu.Unlock()
	search, err := s.Lookup(id)
	if err != nil {
		return false, err
	}
	if !search.Found || search.Item == nil || !search.Complete || search.AuthorityIndeterminate {
		return false, nil
	}
	item := *search.Item
	if item.Revision != expectedRevision {
		return false, nil
	}
	if err := verifyItemFiles(filepath.Join(item.Root, item.ID), expectedHashes); err != nil {
		return false, err
	}
	return true, nil
}

// StagedImport identifies a staged item and the hashes needed to commit or
// reject it later. Staging lives beside the final folder, never inside it.
type StagedImport struct {
	OperationID  string            `json:"operation_id"`
	ItemID       string            `json:"item_id"`
	LibraryID    string            `json:"library_id"`
	Root         string            `json:"root"`
	StagingPath  string            `json:"staging_path"`
	OwnerPath    string            `json:"owner_path"`
	FinalPath    string            `json:"final_path"`
	DocumentHash string            `json:"document_hash"`
	LoreHashes   map[string]string `json:"lore_hashes,omitempty"`
	MediaHash    string            `json:"media_hash,omitempty"`
	Revision     string            `json:"revision"`
	PlanHash     string            `json:"plan_hash"`
}

type importOwner struct {
	Magic        string            `json:"magic"`
	OperationID  string            `json:"operation_id"`
	ItemID       string            `json:"item_id"`
	Category     string            `json:"category"`
	PlanHash     string            `json:"plan_hash"`
	DocumentHash string            `json:"document_hash"`
	LoreHashes   map[string]string `json:"lore_hashes,omitempty"`
	MediaHash    string            `json:"media_hash,omitempty"`
	Revision     string            `json:"revision"`
}

// ValidID reports whether an identity is safe for a canonical item folder.
func ValidID(id string) bool { return validID.MatchString(strings.TrimSpace(id)) }

// ProposeID returns a new opaque identity for this store's category. It is
// never derived from a title or folder name.
func (s *Store) ProposeID() string {
	if s == nil {
		return ""
	}
	return newID(s.category)
}

// StageImport writes one item into operation-scoped staging. It performs no
// canonical write and does not need the process write lock; commit does.
func (s *Store) StageImport(item ImportItem, operationID string) (StagedImport, error) {
	if s == nil || s.reg == nil {
		return StagedImport{}, fmt.Errorf("%w: content library registry is not initialised", ErrUnavailable)
	}
	operationID = strings.TrimSpace(operationID)
	if !validID.MatchString(operationID) {
		return StagedImport{}, fmt.Errorf("%w: operation identity %q", ErrInvalidDocument, operationID)
	}
	id := strings.TrimSpace(item.ID)
	if !validID.MatchString(id) {
		return StagedImport{}, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	if item.Document == nil {
		return StagedImport{}, fmt.Errorf("%w: document is nil", ErrInvalidDocument)
	}
	if inner, ok := item.Document["id"].(string); ok && strings.TrimSpace(inner) != id {
		return StagedImport{}, fmt.Errorf("%w: document identity %q, want %q", ErrInvalidDocument, inner, id)
	}
	item.Document["id"] = id

	view, ok := s.destinationView(item.DestinationLibraryID)
	if !ok {
		return StagedImport{}, ErrNoDestination
	}
	if view.Category != s.category || !view.Exists || !view.Readable || !view.Role.CanSave() || !view.Role.CanWrite() {
		return StagedImport{}, fmt.Errorf("%w: destination %q lacks save+write capability", ErrNoDestination, item.DestinationLibraryID)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	createdAt := strings.TrimSpace(item.CreatedAt)
	if createdAt == "" {
		createdAt = strings.TrimSpace(stringField(item.Document, "createdAt"))
	}
	if createdAt == "" {
		createdAt = strings.TrimSpace(stringField(item.Document, "created_at"))
	}
	if createdAt == "" {
		createdAt = now
	}
	origin := item.Origin
	if strings.TrimSpace(origin.Kind) == "" {
		origin.Kind = "imported"
	}
	doc := Document{
		Format:    s.format,
		Version:   Version,
		ID:        id,
		Revision:  Revision(item.Document),
		CreatedAt: createdAt,
		UpdatedAt: now,
		Origin:    origin,
		Document:  item.Document,
	}
	payload, err := encodeDocument(doc)
	if err != nil {
		return StagedImport{}, err
	}

	lore := make(map[string][]byte, len(item.Lore))
	for filename, data := range item.Lore {
		if filename != filepath.Base(filename) || strings.ContainsAny(filename, `/\\`) || !s.loreFileAllowed(filename) {
			return StagedImport{}, fmt.Errorf("%w: lore file %q", ErrInvalidDocument, filename)
		}
		if len(data) > maxLoreBytes {
			return StagedImport{}, fmt.Errorf("%w: lore file %q exceeds %d bytes", ErrInvalidDocument, filename, maxLoreBytes)
		}
		var parsed struct {
			Entries []json.RawMessage `json:"entries"`
		}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return StagedImport{}, fmt.Errorf("%w: decode lore %q: %v", ErrInvalidDocument, filename, err)
		}
		lore[filename] = append([]byte(nil), data...)
	}
	if len(item.Media) > maxDocumentBytes {
		return StagedImport{}, fmt.Errorf("%w: media.json exceeds %d bytes", ErrInvalidDocument, maxDocumentBytes)
	}
	normalizedMedia, err := NormalizeMediaManifest(item.Media)
	if err != nil {
		return StagedImport{}, fmt.Errorf("%w: media.json: %v", ErrInvalidDocument, err)
	}
	item.Media = normalizedMedia

	stagingRoot := filepath.Join(view.Path, ImportStagingDir, operationID)
	stagingPath := filepath.Join(stagingRoot, id)
	ownerPath := filepath.Join(stagingRoot, id+"."+importOwnerFile)
	if _, err := os.Lstat(stagingPath); err == nil {
		return StagedImport{}, fmt.Errorf("%w: %s", ErrImportStagingExists, stagingPath)
	} else if !os.IsNotExist(err) {
		return StagedImport{}, fmt.Errorf("inspect staging path: %w", err)
	}
	if err := os.MkdirAll(stagingPath, 0700); err != nil {
		return StagedImport{}, fmt.Errorf("create staging item: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(stagingPath)
		_ = os.Remove(ownerPath)
		_ = os.Remove(stagingRoot)
	}

	documentHash := hashBytes(payload)
	loreHashes := make(map[string]string, len(lore))
	for filename, data := range lore {
		loreHashes[filename] = hashBytes(data)
	}
	mediaHash := ""
	if len(item.Media) > 0 {
		mediaHash = hashBytes(item.Media)
	}
	owner := importOwner{
		Magic:        importOwnerMagic,
		OperationID:  operationID,
		ItemID:       id,
		Category:     string(s.category),
		PlanHash:     item.PlanHash,
		DocumentHash: documentHash,
		LoreHashes:   loreHashes,
		MediaHash:    mediaHash,
		Revision:     doc.Revision,
	}
	ownerData, err := json.MarshalIndent(owner, "", "  ")
	if err != nil {
		cleanup()
		return StagedImport{}, fmt.Errorf("encode import owner: %w", err)
	}
	if err := os.WriteFile(ownerPath, ownerData, 0600); err != nil {
		cleanup()
		return StagedImport{}, fmt.Errorf("write import owner: %w", err)
	}
	if err := os.WriteFile(filepath.Join(stagingPath, s.file), payload, 0644); err != nil {
		cleanup()
		return StagedImport{}, fmt.Errorf("write staged document: %w", err)
	}
	for filename, data := range lore {
		if err := os.WriteFile(filepath.Join(stagingPath, filename), data, 0644); err != nil {
			cleanup()
			return StagedImport{}, fmt.Errorf("write staged lore %q: %w", filename, err)
		}
	}
	if len(item.Media) > 0 {
		if err := os.WriteFile(filepath.Join(stagingPath, "media.json"), item.Media, 0644); err != nil {
			cleanup()
			return StagedImport{}, fmt.Errorf("write staged media: %w", err)
		}
	}

	return StagedImport{
		OperationID:  operationID,
		ItemID:       id,
		LibraryID:    view.ID,
		Root:         view.Path,
		StagingPath:  stagingPath,
		OwnerPath:    ownerPath,
		FinalPath:    filepath.Join(view.Path, id),
		DocumentHash: documentHash,
		LoreHashes:   loreHashes,
		MediaHash:    mediaHash,
		Revision:     doc.Revision,
		PlanHash:     item.PlanHash,
	}, nil
}

// CommitImport atomically publishes a staged item only if the destination is
// still authoritative and the final path is absent. It shares writeMu with
// Save and Update, so an import cannot race an ordinary create or update.
func (s *Store) CommitImport(staged StagedImport) (SaveResult, error) {
	res := SaveResult{Status: "write_failed", ID: staged.ItemID}
	if s == nil || s.reg == nil {
		res.Status = "destination_unavailable"
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	if !validID.MatchString(staged.OperationID) || !validID.MatchString(staged.ItemID) {
		return res, fmt.Errorf("%w: invalid staged identity", ErrInvalidDocument)
	}
	if err := s.validateStagedPaths(staged); err != nil {
		res.Status = "invalid_staging"
		res.Message = err.Error()
		return res, nil
	}

	writeMu.Lock()
	defer writeMu.Unlock()

	view, ok := s.destinationView(staged.LibraryID)
	if !ok || view.Category != s.category || filepath.Clean(view.Path) != filepath.Clean(staged.Root) {
		res.Status = "destination_unavailable"
		res.Message = "The destination registration changed or is unavailable; nothing was written."
		return res, nil
	}
	if !view.Exists || !view.Readable || !view.Role.CanSave() || !view.Role.CanWrite() {
		res.Status = "permission_denied"
		res.Message = "The destination no longer grants save and write capability."
		return res, nil
	}
	if err := s.verifyImportOwner(staged); err != nil {
		res.Status = "invalid_staging"
		res.Message = err.Error()
		return res, nil
	}
	if _, err := os.Lstat(staged.FinalPath); err == nil {
		res.Status = "collision"
		res.Message = "An item with this identity already exists."
		return res, nil
	} else if !os.IsNotExist(err) {
		res.Status = "write_failed"
		res.Message = err.Error()
		return res, nil
	}

	search, err := s.Lookup(staged.ItemID)
	if err != nil {
		res.Status = "invalid_document"
		res.Message = err.Error()
		return res, nil
	}
	if search.Found {
		res.Status = "collision"
		res.Message = "An item with this identity already exists in a canonical root."
		return res, nil
	}
	if !search.Complete {
		res.Status = "destination_unavailable"
		res.Message = "No complete search of the item roots could prove this identity is absent; nothing was written."
		return res, nil
	}
	if search.AuthorityIndeterminate {
		res.Status = "authority_indeterminate"
		res.Message = "The destination authority is indeterminate; nothing was written."
		return res, nil
	}
	if err := s.verifyStagedFiles(staged); err != nil {
		res.Status = "invalid_staging"
		res.Message = err.Error()
		return res, nil
	}

	if err := os.Rename(staged.StagingPath, staged.FinalPath); err != nil {
		if _, statErr := os.Lstat(staged.FinalPath); statErr == nil {
			res.Status = "collision"
			res.Message = "An item with this identity appeared during commit."
			return res, nil
		}
		res.Status = "write_failed"
		res.Message = err.Error()
		return res, nil
	}
	_ = os.Remove(staged.OwnerPath)
	_ = os.Remove(filepath.Dir(staged.OwnerPath))
	_ = os.Remove(filepath.Join(staged.Root, ImportStagingDir))

	res.Status = "saved"
	res.Path = filepath.Join(staged.FinalPath, s.file)
	res.Revision = staged.Revision
	return res, nil
}

// AbortImport removes only the staging item and owner for this exact staged
// import. It never touches a final item.
func (s *Store) AbortImport(staged StagedImport) error {
	if s == nil || s.reg == nil {
		return fmt.Errorf("%w: content library registry is not initialised", ErrUnavailable)
	}
	writeMu.Lock()
	defer writeMu.Unlock()
	if err := s.validateStagedPaths(staged); err != nil {
		return err
	}
	if err := s.verifyImportOwner(staged); err != nil {
		return err
	}
	if _, err := os.Lstat(staged.FinalPath); err == nil {
		return fmt.Errorf("%w: final item exists; refusing to abort a committed import", ErrImportConflict)
	}
	if err := s.verifyStagedFiles(staged); err != nil {
		return fmt.Errorf("%w: staged contents are not exact: %v", ErrImportConflict, err)
	}
	if err := os.RemoveAll(staged.StagingPath); err != nil {
		return err
	}
	if err := os.Remove(staged.OwnerPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = os.Remove(filepath.Dir(staged.OwnerPath))
	_ = os.Remove(filepath.Join(staged.Root, ImportStagingDir))
	return nil
}

// RemoveImported removes one journal-created item only if its revision and
// every journaled file hash still match, and the folder contains no unexpected
// files. It is used only by explicit rollback.
func (s *Store) RemoveImported(id, expectedRevision string, expectedHashes map[string]string) (SaveResult, error) {
	id = strings.TrimSpace(id)
	res := SaveResult{Status: "rollback_conflict", ID: id}
	if s == nil || s.reg == nil {
		res.Status = "destination_unavailable"
		return res, nil
	}
	if !validID.MatchString(id) {
		return res, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	if len(expectedHashes) == 0 {
		res.Message = "No journaled file hashes were supplied; refusing to remove an item."
		return res, nil
	}

	writeMu.Lock()
	defer writeMu.Unlock()

	search, err := s.Lookup(id)
	if err != nil {
		res.Message = err.Error()
		return res, nil
	}
	if !search.Found || search.Item == nil {
		if !search.Complete {
			res.Status = "destination_unavailable"
			res.Message = "The item roots cannot prove this item is absent; nothing was removed."
			return res, nil
		}
		res.Status = "not_found"
		res.Message = "The imported item is already absent."
		return res, nil
	}
	if search.AuthorityIndeterminate {
		res.Status = "authority_indeterminate"
		res.Message = "The item authority is indeterminate; nothing was removed."
		return res, nil
	}
	item := *search.Item
	if item.Revision != expectedRevision {
		res.Message = "The item revision changed after import; nothing was removed."
		return res, nil
	}
	view, ok := s.registrationView(item.LibraryID)
	if !ok || !view.Role.CanWrite() || !view.Exists {
		res.Status = "permission_denied"
		res.Message = "The holding root no longer grants write permission."
		return res, nil
	}
	if err := verifyItemFiles(filepath.Join(item.Root, item.ID), expectedHashes); err != nil {
		res.Message = err.Error()
		return res, nil
	}
	if err := os.RemoveAll(filepath.Join(item.Root, item.ID)); err != nil {
		res.Status = "write_failed"
		res.Message = err.Error()
		return res, nil
	}
	res.Status = "rolled_back"
	res.Path = filepath.Join(item.Root, item.ID)
	return res, nil
}

func (s *Store) destinationView(id string) (library.View, bool) {
	return s.registrationView(id)
}

func (s *Store) registrationView(id string) (library.View, bool) {
	for _, view := range s.reg.List(s.category) {
		if view.ID == id {
			return view, true
		}
	}
	return library.View{}, false
}

func (s *Store) validateStagedPaths(staged StagedImport) error {
	root := filepath.Clean(staged.Root)
	wantStaging := filepath.Join(root, ImportStagingDir, staged.OperationID, staged.ItemID)
	wantOwner := filepath.Join(root, ImportStagingDir, staged.OperationID, staged.ItemID+"."+importOwnerFile)
	wantFinal := filepath.Join(root, staged.ItemID)
	if filepath.Clean(staged.StagingPath) != wantStaging || filepath.Clean(staged.OwnerPath) != wantOwner || filepath.Clean(staged.FinalPath) != wantFinal {
		return fmt.Errorf("%w: staged paths are not under the declared root", ErrImportConflict)
	}
	if !library.IsUnderDir(wantStaging, root) || !library.IsUnderDir(wantFinal, root) {
		return fmt.Errorf("%w: staged path escapes its root", ErrImportConflict)
	}
	view, ok := s.registrationView(staged.LibraryID)
	if !ok || view.Category != s.category || filepath.Clean(view.Path) != root {
		return fmt.Errorf("%w: staged library no longer matches its root", ErrImportConflict)
	}
	return nil
}

func (s *Store) verifyImportOwner(staged StagedImport) error {
	data, err := os.ReadFile(staged.OwnerPath)
	if err != nil {
		return fmt.Errorf("%w: read staging ownership: %v", ErrImportConflict, err)
	}
	var owner importOwner
	if err := json.Unmarshal(data, &owner); err != nil {
		return fmt.Errorf("%w: decode staging ownership: %v", ErrImportConflict, err)
	}
	if owner.Magic != importOwnerMagic || owner.OperationID != staged.OperationID || owner.ItemID != staged.ItemID || owner.Category != string(s.category) || owner.PlanHash != staged.PlanHash || owner.DocumentHash != staged.DocumentHash || owner.Revision != staged.Revision {
		return fmt.Errorf("%w: ownership marker does not match operation/item", ErrImportConflict)
	}
	if len(owner.LoreHashes) != len(staged.LoreHashes) {
		return fmt.Errorf("%w: lore hash set changed", ErrImportConflict)
	}
	for filename, hash := range staged.LoreHashes {
		if owner.LoreHashes[filename] != hash {
			return fmt.Errorf("%w: lore hash %q changed", ErrImportConflict, filename)
		}
	}
	if owner.MediaHash != staged.MediaHash {
		return fmt.Errorf("%w: media hash changed", ErrImportConflict)
	}
	return nil
}

func (s *Store) verifyStagedFiles(staged StagedImport) error {
	hashes := map[string]string{s.file: staged.DocumentHash}
	for name, hash := range staged.LoreHashes {
		hashes[name] = hash
	}
	if staged.MediaHash != "" {
		hashes["media.json"] = staged.MediaHash
	}
	return verifyItemFiles(staged.StagingPath, hashes)
}

func verifyItemFiles(dir string, expected map[string]string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read item folder: %w", err)
	}
	if len(entries) != len(expected) {
		return fmt.Errorf("item folder has %d files, journal records %d", len(entries), len(expected))
	}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("unexpected item entry %q", entry.Name())
		}
		want, ok := expected[entry.Name()]
		if !ok {
			return fmt.Errorf("unexpected item file %q", entry.Name())
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read %q: %w", entry.Name(), err)
		}
		if got := hashBytes(data); got != want {
			return fmt.Errorf("item file %q changed", entry.Name())
		}
		seen[entry.Name()] = true
	}
	for name := range expected {
		if !seen[name] {
			return fmt.Errorf("journaled item file %q is missing", name)
		}
	}
	return nil
}

func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
