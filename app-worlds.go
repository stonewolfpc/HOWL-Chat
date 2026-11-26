package main

import (
	"fmt"
	"strings"

	"howl-chat/internal/domaincutover"
	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

// World and scenario documents are canonical in internal/worldstore. These
// bindings are the only renderer write path: they return the document plus the
// revision the next save must cite, and they never write a cache.
func (a *App) worldStore() (*worldstore.Store, error) {
	if a.registries.Libraries() == nil {
		return nil, fmt.Errorf("content library registry not initialised")
	}
	return worldstore.New(a.registries.Libraries(), library.Worlds)
}

func (a *App) scenarioStore() (*worldstore.Store, error) {
	if a.registries.Libraries() == nil {
		return nil, fmt.Errorf("content library registry not initialised")
	}
	return worldstore.New(a.registries.Libraries(), library.Scenarios)
}

// listDocuments is a pure read: it never creates a directory or registers a
// root. An uninitialized or unreadable library yields an empty list, and the
// state is reported through GetContentReadiness instead of being disguised as a
// backend error or an empty library.
func (a *App) listDocuments(store *worldstore.Store, category library.Category) ([]map[string]interface{}, error) {
	items, err := store.List()
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		document := item.Document
		document["revision"] = item.Revision
		out = append(out, document)
	}
	return out, nil
}

func (a *App) GetWorlds() ([]map[string]interface{}, error) {
	store, err := a.worldStore()
	if err != nil {
		return nil, err
	}
	items, err := a.listDocuments(store, library.Worlds)
	if err != nil {
		return nil, fmt.Errorf("GetWorlds: %w", err)
	}
	return items, nil
}

func (a *App) GetScenarios() ([]map[string]interface{}, error) {
	store, err := a.scenarioStore()
	if err != nil {
		return nil, err
	}
	items, err := a.listDocuments(store, library.Scenarios)
	if err != nil {
		return nil, fmt.Errorf("GetScenarios: %w", err)
	}
	return items, nil
}

func (a *App) saveDocument(store *worldstore.Store, category library.Category, document map[string]interface{}, label string) (map[string]interface{}, error) {
	if document == nil {
		return nil, fmt.Errorf("%s: document is required", label)
	}
	id, _ := document["id"].(string)
	id = strings.TrimSpace(id)
	base, _ := document["revision"].(string)
	if id == "" {
		// Only creation goes to the default save root, so only creation is gated
		// on it. Updates must reach the library that owns the document; refusing
		// them here would let a read-only default root disable edits to healthy
		// sibling libraries.
		if err := a.requireCreateWritable(category, label); err != nil {
			return nil, err
		}
		result, err := store.Save(document)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return a.reloadDocument(store, result, label)
	}

	search, err := store.Lookup(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}
	if search.Found && search.Item != nil {
		result, err := store.Update(id, document, base)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return a.reloadDocument(store, result, label)
	}
	// A proven absence is a new document. An incomplete search is not: creating
	// it could collide with a copy in an unavailable root.
	if !search.Complete {
		return nil, fmt.Errorf("%s: storage is unavailable or incompletely searched; nothing was written", label)
	}
	if err := a.requireCreateWritable(category, label); err != nil {
		return nil, err
	}
	delete(document, "revision")
	result, err := store.Save(document)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}
	return a.reloadDocument(store, result, label)
}

func (a *App) reloadDocument(store *worldstore.Store, result worldstore.SaveResult, label string) (map[string]interface{}, error) {
	if result.Status != "saved" && result.Status != "updated" {
		message := result.Message
		if message == "" {
			message = result.Status
		}
		return nil, fmt.Errorf("%s: %s", label, message)
	}
	item, err := store.Load(result.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: reload canonical document: %w", label, err)
	}
	document := item.Document
	document["revision"] = item.Revision
	return document, nil
}

func (a *App) SaveWorld(world map[string]interface{}) (map[string]interface{}, error) {
	store, err := a.worldStore()
	if err != nil {
		return nil, err
	}
	return a.saveDocument(store, library.Worlds, world, "SaveWorld")
}

func (a *App) SaveScenario(scenario map[string]interface{}) (map[string]interface{}, error) {
	store, err := a.scenarioStore()
	if err != nil {
		return nil, err
	}
	return a.saveDocument(store, library.Scenarios, scenario, "SaveScenario")
}

func (a *App) deleteDocument(store *worldstore.Store, id, label string) error {
	result, err := store.Delete(id)
	if err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	if result.Status != "deleted" {
		message := result.Message
		if message == "" {
			message = result.Status
		}
		return fmt.Errorf("%s: %s", label, message)
	}
	return nil
}

func (a *App) importDocuments(store *worldstore.Store, category library.Category, documents []interface{}, label string) (int, []domaincutover.Item, error) {
	if err := a.requireCreateWritable(category, label); err != nil {
		return 0, nil, err
	}
	imported := 0
	items := []domaincutover.Item{}
	for _, raw := range documents {
		document, ok := raw.(map[string]interface{})
		if !ok || document == nil {
			continue
		}
		id, _ := document["id"].(string)
		if strings.TrimSpace(id) == "" {
			continue
		}
		search, err := store.Lookup(id)
		if err != nil {
			return imported, items, fmt.Errorf("%s: %w", label, err)
		}
		if !search.Complete {
			return imported, items, fmt.Errorf("%s: storage is unavailable or incompletely searched; nothing was imported", label)
		}
		// Canonical content wins. A legacy item only seeds an absent identity, so
		// a stale renderer array can never overwrite a newer canonical save.
		if search.Found {
			continue
		}
		delete(document, "revision")
		result, err := store.Save(document)
		if err != nil {
			return imported, items, fmt.Errorf("%s: %w", label, err)
		}
		if result.Status != "saved" {
			return imported, items, fmt.Errorf("%s: %s", label, result.Message)
		}
		hashes, _ := store.ItemHashes(result.ID)
		items = append(items, domaincutover.Item{
			ID:       result.ID,
			Action:   "created",
			Revision: result.Revision,
			Path:     result.Path,
			Hashes:   hashes,
		})
		imported++
	}
	return imported, items, nil
}

// ImportWorlds seeds canonical storage from the renderer's pre-cutover array.
// It never replaces existing canonical content.
func (a *App) ImportWorlds(worlds []interface{}) (int, error) {
	store, err := a.worldStore()
	if err != nil {
		return 0, err
	}
	imported, items, err := a.importDocuments(store, library.Worlds, worlds, "ImportWorlds")
	if err == nil || len(items) > 0 {
		if recordErr := a.recordWorldCutover("world", imported, items); recordErr != nil && err == nil {
			err = recordErr
		}
	}
	return imported, err
}

func (a *App) ImportScenarios(scenarios []interface{}) (int, error) {
	store, err := a.scenarioStore()
	if err != nil {
		return 0, err
	}
	imported, items, err := a.importDocuments(store, library.Scenarios, scenarios, "ImportScenarios")
	if err == nil || len(items) > 0 {
		if recordErr := a.recordWorldCutover("scenario", imported, items); recordErr != nil && err == nil {
			err = recordErr
		}
	}
	return imported, err
}

func (a *App) DeleteWorld(id string) error {
	store, err := a.worldStore()
	if err != nil {
		return err
	}
	return a.deleteDocument(store, id, "DeleteWorld")
}

func (a *App) DeleteScenario(id string) error {
	store, err := a.scenarioStore()
	if err != nil {
		return err
	}
	return a.deleteDocument(store, id, "DeleteScenario")
}
