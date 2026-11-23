// Package characterstore owns canonical character documents. A character is
// one folder under the Characters state root, addressed by an opaque stable
// identity: <root>/<id>/character.json.
//
// The store is the only writer of character documents. It reads legacy raw
// character.json files (the pre-envelope layout) and writes howl-character-v1
// envelopes from then on, so existing folders keep working without a migration
// pass. Import is explicit and idempotent; Rollback restores the documents an
// import replaced.
package characterstore

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
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
)

const (
	// Format is the canonical envelope format.
	Format = "howl-character-v1"
	// Version is the envelope schema version.
	Version = 1
	// fileName is the document name inside a character folder.
	fileName = "character.json"
	// maxDocumentBytes bounds one document file.
	maxDocumentBytes = 16 << 20
)

var (
	// ErrInvalidDocument reports a document that cannot be stored safely.
	ErrInvalidDocument = errors.New("character document is invalid")
	// ErrNotFound reports an identity the store does not hold.
	ErrNotFound = errors.New("character is not in the character store")
	// ErrRevisionConflict reports a save based on a stale revision.
	ErrRevisionConflict = errors.New("character changed since it was loaded")
)

// safeID is the identity shape for newly created folders: opaque, path-safe,
// bounded. Legacy folders whose names do not match are still readable and
// updatable in place; only new identities must match.
var safeID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$`)

// Item is one canonical document.
type Item struct {
	ID       string                 `json:"id"`
	Revision string                 `json:"revision"`
	Path     string                 `json:"path"`
	Legacy   bool                   `json:"legacy"`
	Document map[string]interface{} `json:"document"`
}

// SaveResult reports the outcome of one save.
type SaveResult struct {
	Status   string `json:"status"`
	ID       string `json:"id"`
	Revision string `json:"revision"`
	Path     string `json:"path"`
	Message  string `json:"message,omitempty"`
}

// ImportResult reports one imported document and what it replaced, so a caller
// can undo the import deliberately.
type ImportResult struct {
	ID       string                 `json:"id"`
	Status   string                 `json:"status"`
	Revision string                 `json:"revision"`
	Previous map[string]interface{} `json:"previous,omitempty"`
}

type envelope struct {
	Format    string                 `json:"format"`
	Version   int                    `json:"version"`
	ID        string                 `json:"id"`
	Revision  string                 `json:"revision"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Origin    map[string]interface{} `json:"origin,omitempty"`
	Document  map[string]interface{} `json:"document"`
}

// Store is the single writer for the Characters directory.
type Store struct {
	root string
	mu   sync.Mutex
}

// New builds a store rooted at dir.
func New(dir string) *Store {
	return &Store{root: dir}
}

// Root returns the directory the store owns.
func (s *Store) Root() string {
	if s == nil {
		return ""
	}
	return s.root
}

// List returns every readable character. Files that are not valid JSON or do
// not carry a usable identity are skipped rather than failing the whole list,
// matching the old GetCharacters behaviour. Legacy raw documents are returned
// with Legacy=true.
func (s *Store) List() ([]Item, error) {
	if s == nil || s.root == "" {
		return nil, fmt.Errorf("%w: character store is not initialised", ErrInvalidDocument)
	}
	entries, err := os.ReadDir(s.root)
	if err != nil {
		if os.IsNotExist(err) {
			return []Item{}, nil
		}
		return nil, err
	}
	items := make([]Item, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		item, err := s.readFolder(entry.Name())
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

// Get returns one character by its stable identity. The folder name is accepted
// as an identity too, so the old name-based callers keep working.
func (s *Store) Get(id string) (Item, error) {
	if s == nil || s.root == "" {
		return Item{}, fmt.Errorf("%w: character store is not initialised", ErrInvalidDocument)
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return Item{}, ErrNotFound
	}
	entries, err := os.ReadDir(s.root)
	if err != nil {
		if os.IsNotExist(err) {
			return Item{}, ErrNotFound
		}
		return Item{}, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		item, err := s.readFolder(entry.Name())
		if err != nil {
			continue
		}
		if item.ID == id || entry.Name() == id {
			return item, nil
		}
	}
	return Item{}, ErrNotFound
}

// Save creates or updates one character. id is taken from document["id"]; a
// missing id gets a new opaque one. An update accepts the caller's revision in
// baseRevision: when supplied it must match the stored revision, so a stale
// editor cannot overwrite a newer save.
func (s *Store) Save(document map[string]interface{}, baseRevision string) (Item, error) {
	if s == nil || s.root == "" {
		return Item{}, fmt.Errorf("%w: character store is not initialised", ErrInvalidDocument)
	}
	if document == nil {
		return Item{}, fmt.Errorf("%w: document is nil", ErrInvalidDocument)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := strings.TrimSpace(stringField(document, "id"))
	if id == "" {
		id = newID()
	}
	if strings.ContainsAny(id, `/\`) || id == "." || id == ".." {
		return Item{}, fmt.Errorf("%w: identity %q", ErrInvalidDocument, id)
	}

	existing, err := s.Get(id)
	hasExisting := err == nil
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Item{}, err
	}

	dirName := id
	createdAt := strings.TrimSpace(stringField(document, "createdAt"))
	if createdAt == "" {
		createdAt = strings.TrimSpace(stringField(document, "created_at"))
	}
	if hasExisting {
		// Reuse the folder the identity already lives in, even when its name is
		// not a safe new-identity shape (legacy folders may contain spaces).
		dirName = filepath.Base(filepath.Dir(existing.Path))
		if baseRevision = strings.TrimSpace(baseRevision); baseRevision != "" && baseRevision != existing.Revision {
			return Item{}, ErrRevisionConflict
		}
		if createdAt == "" {
			createdAt = strings.TrimSpace(stringField(existing.Document, "createdAt"))
		}
	} else if !safeID.MatchString(id) {
		return Item{}, fmt.Errorf("%w: identity %q is not path-safe", ErrInvalidDocument, id)
	}
	if createdAt == "" {
		createdAt = time.Now().UTC().Format(time.RFC3339)
	}

	merged := document
	if hasExisting {
		merged = mergeDocument(existing.Document, document)
	}
	merged["id"] = id
	delete(merged, "revision")
	delete(merged, "format")
	delete(merged, "version")
	merged["createdAt"] = createdAt

	now := time.Now().UTC().Format(time.RFC3339)
	env := envelope{
		Format:    Format,
		Version:   Version,
		ID:        id,
		Revision:  Revision(merged),
		CreatedAt: createdAt,
		UpdatedAt: now,
		Origin:    map[string]interface{}{"kind": "howl"},
		Document:  merged,
	}
	payload, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return Item{}, fmt.Errorf("%w: encode: %v", ErrInvalidDocument, err)
	}
	if len(payload) > maxDocumentBytes {
		return Item{}, fmt.Errorf("%w: document exceeds %d bytes", ErrInvalidDocument, maxDocumentBytes)
	}

	dir := filepath.Join(s.root, dirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Item{}, err
	}
	target := filepath.Join(dir, fileName)
	if err := os.WriteFile(target, payload, 0644); err != nil {
		return Item{}, err
	}
	return Item{ID: id, Revision: env.Revision, Path: target, Document: merged}, nil
}

// Delete removes one character folder. The folder must be one the store
// discovered under its root; callers cannot supply a path.
func (s *Store) Delete(id string) error {
	if s == nil || s.root == "" {
		return fmt.Errorf("%w: character store is not initialised", ErrInvalidDocument)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	item, err := s.Get(id)
	if err != nil {
		return err
	}
	dir := filepath.Dir(item.Path)
	rel, err := filepath.Rel(s.root, dir)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%w: character path is outside the store", ErrInvalidDocument)
	}
	return os.RemoveAll(dir)
}

// Import upserts each document by identity and records what it replaced. It is
// idempotent: importing the same list twice produces the same documents.
func (s *Store) Import(documents []map[string]interface{}) ([]ImportResult, error) {
	results := make([]ImportResult, 0, len(documents))
	var firstErr error
	for _, document := range documents {
		if document == nil {
			continue
		}
		id := strings.TrimSpace(stringField(document, "id"))
		result := ImportResult{ID: id, Status: "failed"}
		if id != "" {
			if previous, err := s.Get(id); err == nil {
				result.Previous = previous.Document
				result.Status = "updated"
			} else if errors.Is(err, ErrNotFound) {
				result.Status = "created"
			}
		} else {
			result.Status = "created"
		}
		item, err := s.Save(document, "")
		if err != nil {
			result.Status = "failed"
			if firstErr == nil {
				firstErr = err
			}
			results = append(results, result)
			continue
		}
		result.ID = item.ID
		result.Revision = item.Revision
		results = append(results, result)
	}
	return results, firstErr
}

// Rollback restores the documents a prior Import replaced, and deletes the
// characters that import created. It is the deliberate undo path for Import.
func (s *Store) Rollback(results []ImportResult) error {
	var firstErr error
	for _, result := range results {
		if result.ID == "" {
			continue
		}
		if result.Previous != nil {
			if _, err := s.Save(result.Previous, ""); err != nil && firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := s.Delete(result.ID); err != nil && !errors.Is(err, ErrNotFound) && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Revision is the content hash stored in the envelope.
func Revision(document map[string]interface{}) string {
	payload, err := json.Marshal(document)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func (s *Store) readFolder(dirName string) (Item, error) {
	target := filepath.Join(s.root, dirName, fileName)
	payload, err := os.ReadFile(target)
	if err != nil {
		return Item{}, err
	}
	if len(payload) > maxDocumentBytes {
		return Item{}, fmt.Errorf("%w: document exceeds %d bytes", ErrInvalidDocument, maxDocumentBytes)
	}
	var probe struct {
		Format string `json:"format"`
	}
	_ = json.Unmarshal(payload, &probe)
	if probe.Format == Format {
		var env envelope
		dec := json.NewDecoder(bytes.NewReader(payload))
		dec.UseNumber()
		if err := dec.Decode(&env); err != nil {
			return Item{}, err
		}
		if env.Version != Version || env.Document == nil {
			return Item{}, fmt.Errorf("%w: envelope version %d", ErrInvalidDocument, env.Version)
		}
		id := strings.TrimSpace(env.ID)
		if id == "" {
			id = dirName
		}
		revision := env.Revision
		if revision == "" {
			revision = Revision(env.Document)
		}
		return Item{ID: id, Revision: revision, Path: target, Document: env.Document}, nil
	}
	var document map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	if err := dec.Decode(&document); err != nil {
		return Item{}, err
	}
	if document == nil {
		return Item{}, fmt.Errorf("%w: document is empty", ErrInvalidDocument)
	}
	id := strings.TrimSpace(stringField(document, "id"))
	if id == "" {
		id = dirName
	}
	return Item{ID: id, Revision: Revision(document), Path: target, Legacy: true, Document: document}, nil
}

func mergeDocument(existing, patch map[string]interface{}) map[string]interface{} {
	merged := deepCopyMap(existing)
	for key, value := range patch {
		if key == "revision" || key == "format" || key == "version" {
			continue
		}
		if value == nil {
			delete(merged, key)
			continue
		}
		if patchMap, ok := value.(map[string]interface{}); ok {
			if existingMap, ok := merged[key].(map[string]interface{}); ok {
				merged[key] = mergeDocument(existingMap, patchMap)
				continue
			}
		}
		merged[key] = deepCopy(value)
	}
	return merged
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

func newID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("chr_%d", time.Now().UnixNano())
	}
	return "chr_" + hex.EncodeToString(buf)
}
