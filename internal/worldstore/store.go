// Package worldstore owns canonical world and scenario documents. A world or
// scenario is one folder inside a registered content library, addressed by an
// opaque stable identity: <root>/<id>/world.json or <root>/<id>/scenario.json.
// The content-library registry supplies address, roles, and the default save
// destination. This store is the only writer of world/scenario documents.
//
// Import, chat migration, snapshot refresh, and media ownership are deliberately
// absent; the contracts for those live in the frozen F3-C record and are wired
// in later, on top of this foundation.
package worldstore

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/library"
)

const (
	// FormatWorld and FormatScenario are the canonical envelope formats.
	FormatWorld    = "howl-world-v1"
	FormatScenario = "howl-scenario-v1"
	// Version is the envelope schema version.
	Version = 1
	// maxDocumentBytes bounds one document file.
	maxDocumentBytes = 8 << 20
)

var (
	// ErrInvalidDocument reports a document that cannot be stored safely.
	ErrInvalidDocument = errors.New("world/scenario document is invalid")
	// ErrNotFound reports an identity that no readable root holds.
	ErrNotFound = errors.New("world/scenario is not in a canonical library")
	// ErrNoDestination reports that no save root is configured.
	ErrNoDestination = errors.New("no world/scenario save destination is configured")
	// ErrUnavailable reports that no complete search could be made, so absence
	// cannot be proven. It is never equivalent to ErrNotFound.
	ErrUnavailable = errors.New("world/scenario storage is unavailable or incompletely searched")
)

// validID is the frozen identity shape: opaque, path-safe, bounded.
var validID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$`)

// Origin records where an item came from. LegacyIDs keeps reassigned identities
// so references can be migrated deliberately rather than by guessing.
type Origin struct {
	Kind       string   `json:"kind"`
	LegacyIDs  []string `json:"legacy_ids,omitempty"`
	LegacyName string   `json:"legacy_name,omitempty"`
}

// Document is the on-disk envelope. Document stays an open map so page fields
// the store does not understand survive a round trip.
type Document struct {
	Format    string                 `json:"format"`
	Version   int                    `json:"version"`
	ID        string                 `json:"id"`
	Revision  string                 `json:"revision"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Origin    Origin                 `json:"origin"`
	Document  map[string]interface{} `json:"document"`
}

// Item is one canonical document as seen through a registered root.
type Item struct {
	ID               string `json:"id"`
	Format           string `json:"format"`
	Revision         string `json:"revision"`
	ComputedRevision string `json:"computed_revision"`
	RevisionMismatch bool   `json:"revision_mismatch"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	Origin           Origin `json:"origin"`
	LibraryID        string `json:"library_id"`
	Root             string `json:"root"`
	Path             string `json:"path"`
	Available        bool   `json:"available"`
	Editable         bool   `json:"editable"`
	// AuthorityIndeterminate means another root is unavailable or another copy
	// exists, so this copy may be stale and must not be edited or treated as canon.
	AuthorityIndeterminate bool                   `json:"authority_indeterminate,omitempty"`
	Document               map[string]interface{} `json:"document"`
}

// SaveResult is a structured outcome. Callers never infer success from an error
// and never present a cache write as a canonical save.
type SaveResult struct {
	Status          string `json:"status"`
	ID              string `json:"id"`
	Path            string `json:"path,omitempty"`
	Revision        string `json:"revision,omitempty"`
	CurrentRevision string `json:"current_revision,omitempty"`
	Message         string `json:"message,omitempty"`
}

// writeMu serializes read-check-write across every Store in this process, so two
// updates cannot both pass the same base revision and overwrite one another.
// Cross-process writers are outside the single-process desktop runtime's contract;
// the registry's atomic replace still guarantees a complete file on disk.
var writeMu sync.Mutex

// Store is the single writer for one content category.
type Store struct {
	reg      *library.Registry
	category library.Category
	format   string
	file     string
}

// New builds a store for worlds or scenarios. Other categories are rejected so
// this package cannot silently take ownership of content it does not own.
func New(reg *library.Registry, category library.Category) (*Store, error) {
	switch category {
	case library.Worlds:
		return &Store{reg: reg, category: category, format: FormatWorld, file: "world.json"}, nil
	case library.Scenarios:
		return &Store{reg: reg, category: category, format: FormatScenario, file: "scenario.json"}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported category %q", ErrInvalidDocument, category)
	}
}

// Save creates one canonical document. Creating requires a save root; it never
// replaces an existing item.
func (s *Store) Save(document map[string]interface{}) (SaveResult, error) {
	res := SaveResult{Status: "write_failed"}
	if s == nil || s.reg == nil {
		res.Status = "destination_unavailable"
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	if document == nil {
		return res, fmt.Errorf("%w: document is nil", ErrInvalidDocument)
	}

	writeMu.Lock()
	defer writeMu.Unlock()

	id := strings.TrimSpace(stringField(document, "id"))
	if id == "" {
		id = newID(s.category)
		document["id"] = id
	}
	if !validID.MatchString(id) {
		return res, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	res.ID = id

	reg, ok := s.reg.DefaultSave(s.category)
	if !ok {
		return res, ErrNoDestination
	}

	now := time.Now().UTC().Format(time.RFC3339)
	createdAt := strings.TrimSpace(stringField(document, "createdAt"))
	if createdAt == "" {
		createdAt = strings.TrimSpace(stringField(document, "created_at"))
	}
	if createdAt == "" {
		createdAt = now
	}
	doc := Document{
		Format:    s.format,
		Version:   Version,
		ID:        id,
		Revision:  Revision(document),
		CreatedAt: createdAt,
		UpdatedAt: now,
		Origin:    Origin{Kind: "howl"},
		Document:  document,
	}
	payload, err := encodeDocument(doc)
	if err != nil {
		return res, err
	}

	target, err := s.reg.WriteNewFile(reg.ID, filepath.Join(id, s.file), payload)
	if err != nil {
		return res, s.mapWriteError(&res, err)
	}

	res.Status = "saved"
	res.Path = target
	res.Revision = doc.Revision
	return res, nil
}

// Update replaces known fields on an existing document while preserving fields
// the caller does not understand. The caller must supply the revision it loaded;
// a mismatch returns revision_conflict and changes nothing.
func (s *Store) Update(id string, patch map[string]interface{}, baseRevision string) (SaveResult, error) {
	id = strings.TrimSpace(id)
	res := SaveResult{Status: "write_failed", ID: id}
	if s == nil || s.reg == nil {
		res.Status = "destination_unavailable"
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	if !validID.MatchString(id) {
		return res, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}
	if patch == nil {
		return res, fmt.Errorf("%w: patch is nil", ErrInvalidDocument)
	}

	writeMu.Lock()
	defer writeMu.Unlock()

	search, err := s.Lookup(id)
	if err != nil {
		return res, err
	}
	if !search.Found || search.Item == nil {
		if !search.Complete {
			res.Status = "destination_unavailable"
			res.Message = "No readable library could prove whether this item exists; nothing was changed."
			return res, nil
		}
		res.Status = "not_found"
		res.Message = "The world/scenario is not in a canonical library."
		return res, nil
	}
	item := *search.Item
	res.CurrentRevision = item.Revision

	if search.AuthorityIndeterminate {
		res.Status = "authority_indeterminate"
		res.Message = "The item was found, but another root is unavailable or holds another copy; its authority cannot be established, so nothing was changed."
		return res, nil
	}

	if !item.Editable {
		res.Status = "permission_denied"
		res.Message = "The library holding this item does not grant write permission."
		return res, nil
	}
	baseRevision = strings.TrimSpace(baseRevision)
	if baseRevision == "" {
		res.Status = "revision_conflict"
		res.Message = "A base revision is required to update an item."
		return res, nil
	}
	if baseRevision != item.Revision {
		res.Status = "revision_conflict"
		res.Message = "The document changed since it was loaded. Reload before saving."
		return res, nil
	}

	merged, err := mergeDocument(item.Document, patch, id, item.CreatedAt)
	if err != nil {
		return res, err
	}

	doc := Document{
		Format:    s.format,
		Version:   Version,
		ID:        id,
		Revision:  Revision(merged),
		CreatedAt: item.CreatedAt,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Origin:    item.Origin,
		Document:  merged,
	}
	payload, err := encodeDocument(doc)
	if err != nil {
		return res, err
	}

	target, err := s.reg.WriteFile(item.LibraryID, filepath.Join(id, s.file), payload)
	if err != nil {
		return res, s.mapWriteError(&res, err)
	}

	res.Status = "updated"
	res.Path = target
	res.Revision = doc.Revision
	return res, nil
}

// Search is the outcome of looking for one identity. Complete reports that
// every registered root was readable and searched without error; only then may
// absence be asserted. AuthorityIndeterminate reports that a found copy cannot
// be treated as canonical because another root is unavailable or another copy
// exists. Consumers must never treat incompleteness or indeterminate authority
// as permission to recreate, import, or reassign references.
type Search struct {
	Item                   *Item    `json:"item,omitempty"`
	Found                  bool     `json:"found"`
	Complete               bool     `json:"complete"`
	AuthorityIndeterminate bool     `json:"authority_indeterminate,omitempty"`
	Copies                 int      `json:"copies,omitempty"`
	UnavailableRoots       []string `json:"unavailable_roots,omitempty"`
	Errors                 []string `json:"errors,omitempty"`
}

// Lookup searches every registered root. Roots are equivalent search locations,
// not ranked authorities: the registry sorts them by path for reproducible
// retrieval, and the first readable copy is returned for display. A copy is
// authoritative only when every root was searched and exactly one copy exists.
func (s *Store) Lookup(id string) (Search, error) {
	if s == nil || s.reg == nil {
		return Search{}, fmt.Errorf("content library registry is not initialised")
	}
	id = strings.TrimSpace(id)
	if !validID.MatchString(id) {
		return Search{}, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}

	views := s.reg.List(s.category)
	search := Search{Complete: len(views) > 0}
	if len(views) == 0 {
		search.Errors = append(search.Errors, "no world/scenario library is registered")
	}

	candidates := make([]Item, 0, 1)
	for _, view := range views {
		if !view.Exists || !view.Readable {
			search.Complete = false
			search.UnavailableRoots = append(search.UnavailableRoots, view.Path)
			continue
		}
		path := filepath.Join(view.Path, id, s.file)
		payload, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			search.Complete = false
			search.Errors = append(search.Errors, fmt.Sprintf("%s: %v", path, err))
			continue
		}
		doc, err := decodeDocument(payload, s.format, id)
		if err != nil {
			return Search{}, err
		}
		computed := Revision(doc.Document)
		candidates = append(candidates, Item{
			ID:               doc.ID,
			Format:           doc.Format,
			Revision:         doc.Revision,
			ComputedRevision: computed,
			RevisionMismatch: computed != doc.Revision,
			CreatedAt:        doc.CreatedAt,
			UpdatedAt:        doc.UpdatedAt,
			Origin:           doc.Origin,
			LibraryID:        view.ID,
			Root:             view.Path,
			Path:             path,
			Available:        true,
			Editable:         view.Writable,
			Document:         doc.Document,
		})
	}

	search.Copies = len(candidates)
	if len(candidates) > 0 {
		item := candidates[0]
		search.Found = true
		search.Item = &item
		search.AuthorityIndeterminate = len(search.UnavailableRoots) > 0 || len(search.Errors) > 0 || len(candidates) > 1
		if search.AuthorityIndeterminate {
			search.Item.AuthorityIndeterminate = true
			search.Item.Editable = false
		}
	}
	return search, nil
}

// Load returns one canonical document, or ErrNotFound only when absence was
// proven. An unavailable or incompletely searched root returns ErrUnavailable.
// A stored revision that no longer matches the content is flagged, never trusted.
func (s *Store) Load(id string) (Item, error) {
	search, err := s.Lookup(id)
	if err != nil {
		return Item{}, err
	}
	if search.Found && search.Item != nil {
		return *search.Item, nil
	}
	if !search.Complete {
		reason := "incomplete root search"
		if len(search.UnavailableRoots) > 0 {
			reason = "unavailable roots: " + strings.Join(search.UnavailableRoots, ", ")
		} else if len(search.Errors) > 0 {
			reason = strings.Join(search.Errors, "; ")
		}
		return Item{}, fmt.Errorf("%w: %s", ErrUnavailable, reason)
	}
	return Item{}, ErrNotFound
}

func (s *Store) mapWriteError(res *SaveResult, err error) error {
	switch {
	case errors.Is(err, library.ErrCollision):
		res.Status = "collision"
		res.Message = "An item with this identity already exists."
	case errors.Is(err, library.ErrReadOnly):
		res.Status = "permission_denied"
		res.Message = "The library does not grant write permission."
	case errors.Is(err, library.ErrUnsafeTarget):
		res.Status = "unsafe_target"
		res.Message = "The target is not a regular file."
	case errors.Is(err, library.ErrInvalidPath):
		res.Status = "invalid_path"
		res.Message = "The item path is invalid."
	case errors.Is(err, os.ErrPermission):
		res.Status = "permission_denied"
		res.Message = "The destination denied permission to write."
	default:
		res.Status = "write_failed"
		res.Message = err.Error()
	}
	return nil
}

func decodeDocument(payload []byte, format, id string) (Document, error) {
	if len(payload) > maxDocumentBytes {
		return Document{}, fmt.Errorf("%w: document exceeds %d bytes", ErrInvalidDocument, maxDocumentBytes)
	}
	var doc Document
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	if err := dec.Decode(&doc); err != nil {
		return Document{}, fmt.Errorf("%w: decode: %v", ErrInvalidDocument, err)
	}
	if doc.Format != format {
		return Document{}, fmt.Errorf("%w: format %q, want %q", ErrInvalidDocument, doc.Format, format)
	}
	if doc.Version != Version {
		return Document{}, fmt.Errorf("%w: version %d, want %d", ErrInvalidDocument, doc.Version, Version)
	}
	if doc.ID != id {
		return Document{}, fmt.Errorf("%w: envelope identity %q, want %q", ErrInvalidDocument, doc.ID, id)
	}
	if doc.Document == nil {
		return Document{}, fmt.Errorf("%w: document body is missing", ErrInvalidDocument)
	}
	if inner, ok := doc.Document["id"].(string); ok && strings.TrimSpace(inner) != id {
		return Document{}, fmt.Errorf("%w: document identity %q, want %q", ErrInvalidDocument, inner, id)
	}
	return doc, nil
}

func encodeDocument(doc Document) ([]byte, error) {
	payload, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: encode: %v", ErrInvalidDocument, err)
	}
	if len(payload) > maxDocumentBytes {
		return nil, fmt.Errorf("%w: document exceeds %d bytes", ErrInvalidDocument, maxDocumentBytes)
	}
	return payload, nil
}

// mergeDocument applies the patch without discarding content the caller does not
// know about. Unknown top-level keys survive, and `howl`/`metadata` merge key by
// key so page-specific data is not erased by an ordinary save.
func mergeDocument(existing, patch map[string]interface{}, id, createdAt string) (map[string]interface{}, error) {
	if value, ok := patch["id"]; ok {
		if got, isString := value.(string); !isString || strings.TrimSpace(got) != id {
			return nil, fmt.Errorf("%w: identity is immutable", ErrInvalidDocument)
		}
	}
	for _, key := range []string{"createdAt", "created_at"} {
		if value, ok := patch[key]; ok {
			got, isString := value.(string)
			if !isString || strings.TrimSpace(got) != createdAt {
				return nil, fmt.Errorf("%w: creation time is immutable", ErrInvalidDocument)
			}
		}
	}

	merged := deepCopyMap(existing)
	for key, value := range patch {
		switch key {
		case "id", "createdAt", "created_at", "revision":
			continue
		}
		if key == "howl" || key == "metadata" {
			if patchMap, ok := value.(map[string]interface{}); ok {
				if existingMap, ok := merged[key].(map[string]interface{}); ok {
					merged[key] = mergeKnownFields(existingMap, patchMap)
					continue
				}
			}
			if value == nil {
				if clearableField(key) {
					delete(merged, key)
				}
				continue
			}
			merged[key] = deepCopy(value)
			continue
		}
		if value == nil {
			if clearableField(key) {
				delete(merged, key)
			}
			continue
		}
		merged[key] = deepCopy(value)
	}
	return merged, nil
}

// mergeKnownFields merges one nested map, preserving keys the patch omits.
func mergeKnownFields(existing, patch map[string]interface{}) map[string]interface{} {
	out := deepCopyMap(existing)
	for key, value := range patch {
		if value == nil {
			if clearableField(key) {
				delete(out, key)
			}
			continue
		}
		if patchMap, ok := value.(map[string]interface{}); ok {
			if existingMap, ok := out[key].(map[string]interface{}); ok {
				out[key] = mergeKnownFields(existingMap, patchMap)
				continue
			}
		}
		out[key] = deepCopy(value)
	}
	return out
}

// clearableField is the allowlist for explicit null. A null elsewhere is ignored
// rather than treated as permission to erase unknown data.
func clearableField(key string) bool {
	switch key {
	case "description", "details", "notes", "image", "images", "thumbnail",
		"thumbnailIndex", "world", "worldId", "image_style", "avatar":
		return true
	}
	return false
}

func deepCopyMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for key, value := range in {
		out[key] = deepCopy(value)
	}
	return out
}

func deepCopy(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return deepCopyMap(typed)
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i, item := range typed {
			out[i] = deepCopy(item)
		}
		return out
	default:
		return value
	}
}

func stringField(document map[string]interface{}, key string) string {
	value, ok := document[key]
	if !ok {
		return ""
	}
	text, _ := value.(string)
	return text
}

func newID(category library.Category) string {
	prefix := "wld_"
	if category == library.Scenarios {
		prefix = "scn_"
	}
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return prefix + fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return prefix + hex.EncodeToString(buf)
}
