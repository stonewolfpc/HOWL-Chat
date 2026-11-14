// Package lorebookstore owns canonical lorebook files. A lorebook entry is one
// JSON document inside a registered lorebooks root; the content-library registry
// supplies the address, permissions, and default save destination. This store is
// the only writer of lorebook content. The browser store is a derived cache.
package lorebookstore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/library"
)

const (
	// Format is the canonical document format written by this store.
	Format = "howl-lorebook-v1"
	// Version is the document schema version.
	Version = 1
	// maxEntryBytes bounds one entry file so listing stays cheap.
	maxEntryBytes = 4 << 20
	// maxRootEntries bounds a single root scan.
	maxRootEntries = 5000
)

var (
	// ErrNotCutover is returned when a canonical write is attempted before the
	// explicit import has transferred ownership.
	ErrNotCutover = errors.New("lorebook ownership has not been imported yet")
	// ErrInvalidEntry is returned for an entry that cannot be stored safely.
	ErrInvalidEntry = errors.New("lorebook entry is invalid")
	// ErrNotFound is returned when a canonical entry does not exist.
	ErrNotFound = errors.New("lorebook entry not found")
	// ErrNoDestination is returned when no save root is configured.
	ErrNoDestination = errors.New("no lorebook save destination is configured")
)

var validID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$`)

// Document is the on-disk envelope. The entry itself is kept as an open map so
// page fields survive a round trip without a second schema copy.
type Document struct {
	Format    string                 `json:"format"`
	Version   int                    `json:"version"`
	ID        string                 `json:"id"`
	UpdatedAt string                 `json:"updated_at"`
	Entry     map[string]interface{} `json:"entry"`
}

// Item is one canonical entry as seen through a registered root.
type Item struct {
	ID        string                 `json:"id"`
	LibraryID string                 `json:"library_id"`
	Root      string                 `json:"root"`
	Path      string                 `json:"path"`
	Origin    library.Origin         `json:"origin"`
	Editable  bool                   `json:"editable"`
	Available bool                   `json:"available"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
	Entry     map[string]interface{} `json:"entry"`
}

// RootStatus reports one registered lorebooks root without exposing its files.
type RootStatus struct {
	ID        string         `json:"id"`
	Path      string         `json:"path"`
	Origin    library.Origin `json:"origin"`
	Available bool           `json:"available"`
	Writable  bool           `json:"writable"`
	SaveRoot  bool           `json:"save_root"`
	Count     int            `json:"count"`
}

// ListResult is the complete canonical view plus the recovery-relevant facts.
type ListResult struct {
	Cutover   bool         `json:"cutover"`
	CutoverAt string       `json:"cutover_at,omitempty"`
	Count     int          `json:"count"`
	Items     []Item       `json:"items"`
	Roots     []RootStatus `json:"roots"`
	Conflicts []string     `json:"conflicts"`
	Ignored   int          `json:"ignored"`
	Errors    []string     `json:"errors"`
}

// SaveResult is a structured outcome; callers never infer success from an error.
type SaveResult struct {
	Status  string `json:"status"`
	ID      string `json:"id"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
}

// ImportItem records what happened to one legacy entry.
type ImportItem struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Path   string `json:"path,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// ImportResult reports the explicit ownership transition.
type ImportResult struct {
	Status    string       `json:"status"`
	Cutover   bool         `json:"cutover"`
	RootID    string       `json:"root_id,omitempty"`
	Created   int          `json:"created"`
	Adopted   int          `json:"adopted"`
	Skipped   int          `json:"skipped"`
	Conflicts int          `json:"conflicts"`
	Items     []ImportItem `json:"items"`
	Message   string       `json:"message"`
}

type marker struct {
	Version   int    `json:"version"`
	CutoverAt string `json:"cutover_at"`
	RootID    string `json:"root_id"`
	Count     int    `json:"count"`
}

type journalItem struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Action string `json:"action"`
	Hash   string `json:"hash,omitempty"`
}

type journal struct {
	Version    int           `json:"version"`
	ImportedAt string        `json:"imported_at"`
	RootID     string        `json:"root_id"`
	Status     string        `json:"status"`
	Items      []journalItem `json:"items"`
}

// Store is the single owner of canonical lorebook files.
type Store struct {
	mu               sync.RWMutex
	indexMu          sync.Mutex
	reg              *library.Registry
	markerPath       string
	journalPath      string
	indexPath        string
	indexHash        string
	ownerMarkerPath  string
	ownerJournalPath string
}

// New builds a store over the given registry. A nil registry is allowed; every
// operation then reports that no lorebook library is configured.
func New(reg *library.Registry, configDir string) *Store {
	return &Store{
		reg:              reg,
		markerPath:       filepath.Join(configDir, "lorebook-cutover.json"),
		journalPath:      filepath.Join(configDir, "lorebook-import-journal.json"),
		indexPath:        filepath.Join(configDir, "lorebook-index.json"),
		ownerMarkerPath:  filepath.Join(configDir, "lorebook-owner-migration.json"),
		ownerJournalPath: filepath.Join(configDir, "lorebook-owner-migration-journal.json"),
	}
}

// Cutover reports whether ownership has been explicitly imported.
func (s *Store) Cutover() bool {
	_, err := os.Stat(s.markerPath)
	return err == nil
}

func (s *Store) readMarker() (marker, bool) {
	data, err := os.ReadFile(s.markerPath)
	if err != nil {
		return marker{}, false
	}
	var m marker
	if err := json.Unmarshal(data, &m); err != nil {
		return marker{}, false
	}
	return m, true
}

// List reads every readable lorebooks root once, sequentially. Duplicate IDs
// across roots are reported, not merged. Files that are not this store's format
// are ignored and counted; nothing is ever imported by listing.
func (s *Store) List() ListResult {
	result := ListResult{Items: []Item{}, Roots: []RootStatus{}, Conflicts: []string{}, Errors: []string{}}
	if s.reg == nil {
		result.Errors = append(result.Errors, "content library registry is not initialised")
		return result
	}
	if m, ok := s.readMarker(); ok {
		result.Cutover = true
		result.CutoverAt = m.CutoverAt
	}

	index := s.loadIndex()
	dirty := false
	seen := map[string]string{}
	for _, view := range s.reg.List(library.Lorebooks) {
		root := RootStatus{
			ID:        view.ID,
			Path:      view.Path,
			Origin:    view.Origin,
			Available: view.Readable,
			Writable:  view.Writable,
			SaveRoot:  view.SaveRoot,
		}
		var items []Item
		if view.Readable {
			scanned, ignored, err := s.scanRoot(view)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", view.Path, err))
				root.Available = false
			} else {
				items = scanned
				index.Roots[view.ID] = cachedRoot{
					ScannedAt: time.Now().UTC().Format(time.RFC3339),
					Items:     snapshotItems(items),
				}
				dirty = true
			}
			result.Ignored += ignored
		}
		// A disconnected or unreadable root keeps its last-known entries visible
		// as unavailable; an empty but readable root is genuinely empty.
		if !view.Readable {
			if cached, ok := index.Roots[view.ID]; ok {
				items = unavailableItems(cached.Items)
			}
		}
		root.Count = len(items)
		for _, item := range items {
			if first, dup := seen[item.ID]; dup {
				result.Conflicts = append(result.Conflicts, fmt.Sprintf("%s in both %s and %s", item.ID, first, item.Root))
				continue
			}
			seen[item.ID] = item.Root
			result.Items = append(result.Items, item)
		}
		result.Roots = append(result.Roots, root)
	}
	if dirty {
		s.writeIndex(index)
	}
	sort.SliceStable(result.Items, func(i, j int) bool {
		return strings.ToLower(entryTitle(result.Items[i].Entry)) < strings.ToLower(entryTitle(result.Items[j].Entry))
	})
	result.Count = len(result.Items)
	return result
}

func (s *Store) scanRoot(view library.View) ([]Item, int, error) {
	dirEntries, err := os.ReadDir(view.Path)
	if err != nil {
		return nil, 0, fmt.Errorf("read lorebook root: %w", err)
	}
	items := make([]Item, 0, len(dirEntries))
	ignored := 0
	for i, dirEntry := range dirEntries {
		if i >= maxRootEntries {
			break
		}
		name := dirEntry.Name()
		if dirEntry.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}
		info, err := dirEntry.Info()
		if err != nil || !info.Mode().IsRegular() || info.Size() > maxEntryBytes {
			continue
		}
		path := filepath.Join(view.Path, name)
		entry, updated, ok := readDocument(path)
		if !ok {
			ignored++
			continue
		}
		id := entryID(entry)
		if !validID.MatchString(id) {
			ignored++
			continue
		}
		items = append(items, Item{
			ID:        id,
			LibraryID: view.ID,
			Root:      view.Path,
			Path:      path,
			Origin:    view.Origin,
			Editable:  view.Writable,
			Available: true,
			UpdatedAt: updated,
			Entry:     entry,
		})
	}
	return items, ignored, nil
}

func readDocument(path string) (map[string]interface{}, string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", false
	}
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, "", false
	}
	if doc.Format != Format || doc.Entry == nil {
		return nil, "", false
	}
	entry := doc.Entry
	if entryID(entry) == "" && doc.ID != "" {
		entry["id"] = doc.ID
	}
	if entryID(entry) == "" {
		return nil, "", false
	}
	if entry["updatedAt"] == nil && doc.UpdatedAt != "" {
		entry["updatedAt"] = doc.UpdatedAt
	}
	return entry, doc.UpdatedAt, true
}

// Save writes one canonical entry. Creating requires a save root; replacing
// requires the write role on the root that already holds the item.
func (s *Store) Save(entry map[string]interface{}, replace bool) (SaveResult, error) {
	res := SaveResult{Status: "failed"}
	if s.reg == nil {
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	if entry == nil {
		return res, fmt.Errorf("%w: entry is nil", ErrInvalidEntry)
	}
	if !s.Cutover() {
		return res, ErrNotCutover
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := entryID(entry)
	if id == "" {
		id = newID()
		entry["id"] = id
	}
	if !validID.MatchString(id) {
		return res, fmt.Errorf("%w: identity %q", ErrInvalidEntry, id)
	}
	res.ID = id

	if updated := strings.TrimSpace(stringField(entry, "updatedAt")); updated == "" {
		entry["updatedAt"] = time.Now().UTC().Format(time.RFC3339)
	}
	rel := id + ".json"

	rootID := ""
	if replace {
		item, ok := s.itemByID(id)
		if !ok {
			res.Status = "not_found"
			res.Message = "The entry is not in a canonical library."
			return res, nil
		}
		if !item.Available {
			res.Status = "destination_unavailable"
			res.Message = "The library holding this entry is unavailable."
			return res, nil
		}
		// Ownership is preserved for legacy payloads that do not carry the field
		// at all. A payload that supplies "ownerId" - including an empty string -
		// is explicit: it sets or clears the owner and must not be overridden.
		if _, supplied := entry["ownerId"]; !supplied {
			if existingOwner, ok := item.Entry["ownerId"]; ok {
				entry["ownerId"] = existingOwner
			}
		}
		// ownerName is descriptive only. Preserve it when the payload does not
		// carry it, but an explicit owner clear also clears the stale label.
		if _, supplied := entry["ownerName"]; !supplied {
			if existingName, ok := item.Entry["ownerName"]; ok {
				entry["ownerName"] = existingName
			}
		}
		if owner, supplied := entry["ownerId"]; supplied && owner == "" {
			entry["ownerName"] = ""
		}
		rootID = item.LibraryID
	} else {
		reg, ok := s.reg.DefaultSave(library.Lorebooks)
		if !ok {
			return res, ErrNoDestination
		}
		rootID = reg.ID
	}

	payload, err := encodeDocument(id, entry)
	if err != nil {
		return res, err
	}
	var target string
	if replace {
		target, err = s.reg.WriteFile(rootID, rel, payload)
	} else {
		target, err = s.reg.WriteNewFile(rootID, rel, payload)
	}
	if err != nil {
		switch {
		case errors.Is(err, library.ErrCollision):
			res.Status = "collision"
			res.Message = "An entry with this identity already exists."
		case errors.Is(err, library.ErrReadOnly):
			res.Status = "read_only"
			res.Message = "The selected library does not grant write permission."
		case errors.Is(err, library.ErrUnsafeTarget):
			res.Status = "unsafe_target"
			res.Message = "The target is a directory, symlink, or other non-regular file."
		case errors.Is(err, library.ErrInvalidPath):
			res.Status = "failed"
			res.Message = "The entry path is invalid."
		default:
			res.Status = "failed"
			res.Message = err.Error()
		}
		return res, nil
	}
	res.Path = target
	if replace {
		res.Status = "replaced"
	} else {
		res.Status = "saved"
	}
	return res, nil
}

// Delete removes one canonical entry from its writable root. The path is
// re-derived from the identity and re-checked against the root before removal.
func (s *Store) Delete(id string) (SaveResult, error) {
	res := SaveResult{Status: "failed", ID: id}
	if s.reg == nil {
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	if !validID.MatchString(id) {
		res.Message = "The entry identity is invalid."
		return res, nil
	}
	if !s.Cutover() {
		return res, ErrNotCutover
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.itemByID(id)
	if !ok {
		res.Status = "not_found"
		res.Message = "The entry is not in a canonical library."
		return res, nil
	}
	if !item.Editable {
		res.Status = "read_only"
		res.Message = "This entry belongs to a read-only library."
		return res, nil
	}
	expected := filepath.Join(item.Root, id+".json")
	if filepath.Clean(expected) != filepath.Clean(item.Path) || !library.IsUnderDir(item.Path, item.Root) {
		res.Message = "Refusing to delete an entry outside its registered library."
		return res, nil
	}
	info, err := os.Lstat(item.Path)
	if err != nil {
		res.Status = "not_found"
		res.Message = "The entry file no longer exists."
		return res, nil
	}
	if !info.Mode().IsRegular() {
		res.Status = "unsafe_target"
		res.Message = "The entry file is not a regular file."
		return res, nil
	}
	if err := os.Remove(item.Path); err != nil {
		res.Status = "failed"
		res.Message = err.Error()
		return res, nil
	}
	res.Status = "deleted"
	res.Path = item.Path
	return res, nil
}

// Import is the only ownership transition. It is explicit, create-only,
// idempotent by identity and content hash, and journaled so it can be rolled
// back without touching the legacy browser store.
func (s *Store) Import(entries []map[string]interface{}) (ImportResult, error) {
	res := ImportResult{Status: "failed", Items: []ImportItem{}}
	if s.reg == nil {
		res.Message = "content library registry is not initialised."
		return res, nil
	}
	reg, ok := s.reg.DefaultSave(library.Lorebooks)
	if !ok {
		return res, ErrNoDestination
	}
	res.RootID = reg.ID

	s.mu.Lock()
	defer s.mu.Unlock()

	rec := journal{
		Version:    1,
		ImportedAt: time.Now().UTC().Format(time.RFC3339),
		RootID:     reg.ID,
		Status:     "completed",
		Items:      []journalItem{},
	}
	used := map[string]bool{}
	cutoverAtStart := s.Cutover()

	for _, entry := range entries {
		if entry == nil {
			continue
		}
		action := "created"
		id := entryID(entry)
		if !validID.MatchString(id) || used[id] {
			newIdentity := newID()
			entry["id"] = newIdentity
			action = "reassigned"
			if id != "" {
				res.Items = append(res.Items, ImportItem{ID: id, Action: "reassigned", Detail: "identity was invalid or duplicated; imported as " + newIdentity})
			}
			id = newIdentity
		}
		used[id] = true

		rel := id + ".json"
		target := filepath.Join(reg.Path, rel)
		if !library.IsUnderDir(target, reg.Path) {
			res.Items = append(res.Items, ImportItem{ID: id, Action: "failed", Detail: "path escapes the library"})
			continue
		}
		payload, err := encodeDocument(id, entry)
		if err != nil {
			res.Items = append(res.Items, ImportItem{ID: id, Action: "failed", Detail: err.Error()})
			rec.Status = "partial"
			continue
		}

		if existing, _, ok := readDocument(target); ok {
			if hashEntry(existing) == hashEntry(entry) {
				if cutoverAtStart {
					res.Skipped++
					res.Items = append(res.Items, ImportItem{ID: id, Action: "skipped", Path: target, Detail: "already imported with identical content"})
					continue
				}
				// The file exists while the marker does not, so a previous import
				// was interrupted after writing. Verify it and finish the
				// transition instead of leaving ownership incomplete forever.
				res.Adopted++
				rec.Items = append(rec.Items, journalItem{ID: id, Path: target, Action: "adopted", Hash: hashEntry(entry)})
				res.Items = append(res.Items, ImportItem{ID: id, Action: "adopted", Path: target, Detail: "existing file verified; ownership completed"})
				continue
			}
			res.Conflicts++
			res.Items = append(res.Items, ImportItem{ID: id, Action: "conflict", Path: target, Detail: "a different canonical entry already uses this identity; the existing file was not changed"})
			rec.Status = "partial"
			continue
		}

		if _, err := s.reg.WriteNewFile(reg.ID, rel, payload); err != nil {
			res.Items = append(res.Items, ImportItem{ID: id, Action: "failed", Path: target, Detail: err.Error()})
			rec.Status = "partial"
			continue
		}
		res.Created++
		rec.Items = append(rec.Items, journalItem{ID: id, Path: target, Action: action, Hash: hashEntry(entry)})
		res.Items = append(res.Items, ImportItem{ID: id, Action: action, Path: target})
	}

	if res.Created == 0 && res.Adopted == 0 && res.Conflicts == 0 {
		res.Status = "nothing_to_do"
		if res.Skipped > 0 {
			res.Message = "Every entry is already imported with identical content."
		} else {
			res.Message = "There are no legacy entries to import."
		}
		return res, nil
	}

	if err := s.writeJournal(rec); err != nil {
		return res, err
	}
	if res.Created > 0 || res.Adopted > 0 {
		if err := s.writeMarker(marker{Version: 1, CutoverAt: rec.ImportedAt, RootID: reg.ID, Count: res.Created + res.Adopted}); err != nil {
			return res, err
		}
	}
	res.Cutover = s.Cutover()
	if res.Conflicts > 0 {
		res.Status = "partial"
		res.Message = "Imported entries, but some identities conflicted with existing files and were left untouched."
		return res, nil
	}
	res.Status = "completed"
	if res.Adopted > 0 && res.Created == 0 {
		res.Message = "An interrupted import was verified and ownership completed. The browser copy was left untouched."
	} else {
		res.Message = "Legacy lorebooks are now stored as canonical files. The browser copy was left untouched."
	}
	return res, nil
}

// RollbackImport removes only the files this store created during the last
// import, leaves every other file and the legacy store untouched, and clears
// the cutover marker so the browser store is authoritative again.
func (s *Store) RollbackImport() (ImportResult, error) {
	res := ImportResult{Status: "failed", Items: []ImportItem{}}
	if s.reg == nil {
		res.Message = "content library registry is not initialised."
		return res, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.readJournal()
	if !ok || len(rec.Items) == 0 {
		res.Status = "nothing_to_do"
		res.Message = "No import journal was found; nothing was rolled back."
		return res, nil
	}
	if rec.Status == "rolled_back" {
		res.Status = "nothing_to_do"
		res.Message = "The last import has already been rolled back."
		return res, nil
	}
	reg, ok := s.regView(rec.RootID)
	if !ok {
		res.Message = "The import root is no longer registered; no file was removed."
		return res, nil
	}
	res.RootID = rec.RootID
	for _, item := range rec.Items {
		if item.Action == "adopted" {
			res.Items = append(res.Items, ImportItem{ID: item.ID, Action: "skipped", Detail: "adopted file was not created by this store; left in place"})
			continue
		}
		target := filepath.Join(reg.Path, item.ID+".json")
		if filepath.Clean(target) != filepath.Clean(item.Path) || !library.IsUnderDir(item.Path, reg.Path) {
			res.Items = append(res.Items, ImportItem{ID: item.ID, Action: "failed", Detail: "journal path does not match the registered root"})
			continue
		}
		info, err := os.Lstat(target)
		if err != nil || !info.Mode().IsRegular() {
			res.Items = append(res.Items, ImportItem{ID: item.ID, Action: "skipped", Detail: "file no longer exists"})
			continue
		}
		if err := os.Remove(target); err != nil {
			res.Items = append(res.Items, ImportItem{ID: item.ID, Action: "failed", Detail: err.Error()})
			continue
		}
		res.Items = append(res.Items, ImportItem{ID: item.ID, Action: "removed", Path: target})
		res.Created++
	}
	rec.Status = "rolled_back"
	if err := s.writeJournal(rec); err != nil {
		return res, err
	}
	if err := os.Remove(s.markerPath); err != nil && !os.IsNotExist(err) {
		return res, err
	}
	if err := os.Remove(s.indexPath); err != nil && !os.IsNotExist(err) {
		return res, err
	}
	s.indexMu.Lock()
	s.indexHash = ""
	s.indexMu.Unlock()
	res.Cutover = false
	res.Status = "rolled_back"
	res.Message = "Imported copies were removed. The browser lorebook store was not changed."
	return res, nil
}

func (s *Store) itemByID(id string) (Item, bool) {
	if s.reg == nil {
		return Item{}, false
	}
	for _, view := range s.reg.List(library.Lorebooks) {
		if !view.Readable {
			continue
		}
		items, _, err := s.scanRoot(view)
		if err != nil {
			continue
		}
		for _, item := range items {
			if item.ID == id {
				return item, true
			}
		}
	}
	return Item{}, false
}

func (s *Store) regView(id string) (library.View, bool) {
	if s.reg == nil {
		return library.View{}, false
	}
	for _, view := range s.reg.List(library.Lorebooks) {
		if view.ID == id {
			return view, true
		}
	}
	return library.View{}, false
}

func (s *Store) writeMarker(m marker) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("encode cutover marker: %w", err)
	}
	return writeAtomic(s.markerPath, data)
}

func (s *Store) readJournal() (journal, bool) {
	data, err := os.ReadFile(s.journalPath)
	if err != nil {
		return journal{}, false
	}
	var rec journal
	if err := json.Unmarshal(data, &rec); err != nil {
		return journal{}, false
	}
	return rec, true
}

func (s *Store) writeJournal(rec journal) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("encode import journal: %w", err)
	}
	return writeAtomic(s.journalPath, data)
}

type cachedRoot struct {
	ScannedAt string `json:"scanned_at"`
	Items     []Item `json:"items"`
}

type persistedIndex struct {
	Version int                   `json:"version"`
	Roots   map[string]cachedRoot `json:"roots"`
}

func (s *Store) loadIndex() persistedIndex {
	idx := persistedIndex{Version: 1, Roots: map[string]cachedRoot{}}
	data, err := os.ReadFile(s.indexPath)
	if err != nil {
		return idx
	}
	var stored persistedIndex
	if err := json.Unmarshal(data, &stored); err != nil || stored.Roots == nil {
		return idx
	}
	return stored
}

func (s *Store) writeIndex(idx persistedIndex) {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return
	}
	sum := sha256.Sum256(data)
	key := hex.EncodeToString(sum[:])
	s.indexMu.Lock()
	defer s.indexMu.Unlock()
	if key == s.indexHash {
		return
	}
	if err := writeAtomic(s.indexPath, data); err != nil {
		return
	}
	s.indexHash = key
}

func snapshotItems(items []Item) []Item {
	out := make([]Item, len(items))
	copy(out, items)
	return out
}

func unavailableItems(items []Item) []Item {
	out := make([]Item, len(items))
	for i, item := range items {
		item.Available = false
		item.Editable = false
		out[i] = item
	}
	return out
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create lorebook state directory: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write lorebook state: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("commit lorebook state: %w", err)
	}
	return nil
}

func encodeDocument(id string, entry map[string]interface{}) ([]byte, error) {
	payload, err := json.MarshalIndent(Document{
		Format:    Format,
		Version:   Version,
		ID:        id,
		UpdatedAt: stringField(entry, "updatedAt"),
		Entry:     entry,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: encode: %v", ErrInvalidEntry, err)
	}
	if len(payload) > maxEntryBytes {
		return nil, fmt.Errorf("%w: entry exceeds %d bytes", ErrInvalidEntry, maxEntryBytes)
	}
	return payload, nil
}

func hashEntry(entry map[string]interface{}) string {
	data, err := json.Marshal(entry)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func entryID(entry map[string]interface{}) string {
	return strings.TrimSpace(stringField(entry, "id"))
}

func entryTitle(entry map[string]interface{}) string {
	if title := strings.TrimSpace(stringField(entry, "title")); title != "" {
		return title
	}
	return strings.TrimSpace(stringField(entry, "name"))
}

func stringField(entry map[string]interface{}, key string) string {
	value, ok := entry[key]
	if !ok {
		return ""
	}
	text, _ := value.(string)
	return text
}

func newID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "lb_" + fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return "lb_" + hex.EncodeToString(buf)
}
