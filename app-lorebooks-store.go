package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/library"
	"howl-chat/internal/lorebookstore"
)

// Canonical lorebook bindings. The store owns files; the browser store is a
// derived cache and the legacy source for the explicit import.

func (a *App) lorebookBackend() (*lorebookstore.Store, error) {
	if a.lorebookStore == nil {
		return nil, fmt.Errorf("lorebook store not initialised")
	}
	return a.lorebookStore, nil
}

// GetLorebookStoreStatus reports ownership state without listing entries.
func (a *App) GetLorebookStoreStatus() (lorebookstore.ListResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.ListResult{}, err
	}
	return store.List(), nil
}

// ListLorebooks returns canonical entries from every readable lorebooks root.
func (a *App) ListLorebooks() (lorebookstore.ListResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.ListResult{}, err
	}
	return store.List(), nil
}

// SaveLorebook writes one canonical entry. Creating requires a save root;
// replacing requires the write role on the root that already holds the entry.
func (a *App) SaveLorebook(entry map[string]interface{}, replace bool) (lorebookstore.SaveResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.SaveResult{Status: "failed", Message: err.Error()}, nil
	}
	return store.Save(entry, replace)
}

// DeleteLorebook removes one canonical entry from its writable root.
func (a *App) DeleteLorebook(id string) (lorebookstore.SaveResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.SaveResult{Status: "failed", Message: err.Error()}, nil
	}
	return store.Delete(id)
}

// ImportLorebooks performs the explicit ownership transition from the browser
// store. It is create-only, idempotent, and journaled.
func (a *App) ImportLorebooks(entries []map[string]interface{}) (lorebookstore.ImportResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.ImportResult{Status: "failed", Message: err.Error()}, nil
	}
	return store.Import(entries)
}

// MigrateLorebookOwners rewrites legacy display-name owners to stable character
// IDs supplied by the character library. Ambiguous and unknown names are
// reported and left untouched; successful rewrites are journaled.
func (a *App) MigrateLorebookOwners(identities []lorebookstore.OwnerIdentity, version string) (lorebookstore.OwnerMigrationResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.OwnerMigrationResult{Status: "failed", Message: err.Error()}, nil
	}
	return store.MigrateOwners(identities, version)
}

// RollbackLorebookOwnerMigration restores ownership values recorded by the
// owner migration without touching unrelated entry data.
func (a *App) RollbackLorebookOwnerMigration() (lorebookstore.OwnerMigrationResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.OwnerMigrationResult{Status: "failed", Message: err.Error()}, nil
	}
	return store.RollbackOwners()
}

// RollbackLorebookImport removes only the files created by the last import and
// returns ownership to the browser store.
func (a *App) RollbackLorebookImport() (lorebookstore.ImportResult, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return lorebookstore.ImportResult{Status: "failed", Message: err.Error()}, nil
	}
	return store.RollbackImport()
}

// mergeCanonicalLore overlays canonical file content onto the chat's lorebook
// pool. The chat owns the selection (which books are on, and their enabled
// flags); the canonical store owns the text. Chat sessions persist a snapshot of
// their lorebooks, so without this overlay a chat opened after an edit would
// keep injecting the stale copy. Before cutover, or without a store, the pool is
// returned untouched and the browser store remains the source.
//
// libraryOwned marks a pool that the canonical store owns end-to-end (the
// browser library list). Legacy scoped lore loaded from world/scenario/character
// files is not library-owned and keeps its own lifecycle until it is migrated.
func (a *App) mergeCanonicalLore(entries []lorebook.Entry, libraryOwned bool) []lorebook.Entry {
	if a.lorebookStore == nil || (len(entries) == 0 && !libraryOwned) {
		return entries
	}

	list := a.lorebookStore.List()
	if !list.Cutover {
		return entries
	}

	// List() returns cached last-known items for an unavailable root as well as
	// live ones. Only Available items are current canon; the rest are kept aside
	// so they can be preserved without ever being treated as authoritative.
	canonical := make(map[string]lorebook.Entry, len(list.Items))
	lastKnown := make(map[string]lorebook.Entry, len(list.Items))
	for _, item := range list.Items {
		if item.ID == "" || item.Entry == nil {
			continue
		}
		entry, err := lorebook.ParseLorebookEntryFromMap(item.Entry)
		if err != nil {
			continue
		}
		entry.ID = item.ID
		if item.Available {
			canonical[item.ID] = entry
		} else {
			lastKnown[item.ID] = entry
		}
	}

	// A root that is unavailable or unreadable cannot prove an entry was deleted.
	// Only a complete canonical view may drop snapshot-only IDs.
	complete := canonicalViewComplete(list)

	merged := make([]lorebook.Entry, 0, len(entries))
	for _, entry := range entries {
		current, ok := canonical[entry.ID]
		if !ok {
			// A missing ID in a complete view means the canonical file was deleted;
			// resurrecting it from a chat snapshot would recreate a competing source
			// of truth. In an incomplete view the entry is kept as last-known state,
			// marked stale so it is never injected.
			if libraryOwned {
				if complete {
					continue
				}
				if cached, found := lastKnown[entry.ID]; found {
					cached.Stale = true
					cached.Enabled = entry.Enabled
					merged = append(merged, cached)
					continue
				}
				entry.Stale = true
			}
			merged = append(merged, entry)
			continue
		}
		current.Stale = false
		current.Enabled = entry.Enabled
		merged = append(merged, current)
	}
	return merged
}

func canonicalViewComplete(list lorebookstore.ListResult) bool {
	if len(list.Errors) > 0 {
		return false
	}
	for _, root := range list.Roots {
		if !root.Available {
			return false
		}
	}
	return true
}

// canonicalLorebook returns one canonical entry by ID, if the store is cut over
// and the entry currently exists.
func (a *App) canonicalLorebook(id string) (lorebook.Entry, bool) {
	if a.lorebookStore == nil || id == "" {
		return lorebook.Entry{}, false
	}
	list := a.lorebookStore.List()
	if !list.Cutover {
		return lorebook.Entry{}, false
	}
	for _, item := range list.Items {
		if item.ID != id || item.Entry == nil {
			continue
		}
		entry, err := lorebook.ParseLorebookEntryFromMap(item.Entry)
		if err != nil {
			return lorebook.Entry{}, false
		}
		entry.ID = item.ID
		return entry, true
	}
	return lorebook.Entry{}, false
}

// GetLorebookStoreDiagnostics is a read-only helper for the review surface: it
// reports the journal state without exposing entry content.
func (a *App) GetLorebookStoreDiagnostics() (map[string]interface{}, error) {
	store, err := a.lorebookBackend()
	if err != nil {
		return nil, err
	}
	list := store.List()
	status := map[string]interface{}{
		"cutover":     list.Cutover,
		"cutover_at":  list.CutoverAt,
		"count":       list.Count,
		"conflicts":   list.Conflicts,
		"ignored":     list.Ignored,
		"roots":       list.Roots,
		"errors":      list.Errors,
		"journal":     nil,
		"journal_err": nil,
	}
	configDir := a.getSettingsDir()
	journalPath := filepath.Join(configDir, "lorebook-import-journal.json")
	if data, readErr := os.ReadFile(journalPath); readErr == nil {
		var decoded map[string]interface{}
		if jsonErr := json.Unmarshal(data, &decoded); jsonErr == nil {
			status["journal"] = decoded
		} else {
			status["journal_err"] = jsonErr.Error()
		}
	}
	return status, nil
}

// EnsureLorebookLibrary reports whether a lorebook save root exists. It never
// creates or registers one on the user's behalf.
func (a *App) EnsureLorebookLibrary() (map[string]interface{}, error) {
	out := map[string]interface{}{"configured": false}
	if a.contentLibraries == nil {
		return out, nil
	}
	reg, ok := a.contentLibraries.DefaultSave(library.Lorebooks)
	if !ok {
		return out, nil
	}
	out["configured"] = true
	out["library_id"] = reg.ID
	out["path"] = reg.Path
	return out, nil
}
