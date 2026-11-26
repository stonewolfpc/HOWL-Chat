package main

import (
	"fmt"

	"howl-chat/internal/characterstore"
	"howl-chat/internal/domaincutover"
	"howl-chat/internal/worldstore"
)

// The four legacy domains expose explicit ownership transitions. Reading a page
// never imports or adopts data; the renderer must call one of these methods from
// a user action. Every transition is recorded in configDir so it survives a
// restart and can be undone only from that durable record.

func (a *App) cutoverStore(domain string) *domaincutover.Store {
	return domaincutover.New(a.getSettingsDir(), domain)
}

func (a *App) GetCharacterCutoverStatus() domaincutover.Status {
	return a.cutoverStore("character").Status()
}

func (a *App) GetWorldCutoverStatus() domaincutover.Status {
	return a.cutoverStore("world").Status()
}

func (a *App) GetScenarioCutoverStatus() domaincutover.Status {
	return a.cutoverStore("scenario").Status()
}

func (a *App) GetImageGenerationCutoverStatus() domaincutover.Status {
	return a.cutoverStore("image-generation").Status()
}

// recordCharacterCutover persists the result of an explicit character import.
func (a *App) recordCharacterCutover(results []characterstore.ImportResult) error {
	items := make([]domaincutover.Item, 0, len(results))
	for _, result := range results {
		if result.ID == "" || (result.Status != "created" && result.Status != "updated") {
			continue
		}
		items = append(items, domaincutover.Item{
			ID:       result.ID,
			Action:   result.Status,
			Revision: result.Revision,
			Previous: result.Previous,
		})
	}
	return a.cutoverStore("character").Record(len(items), items, "completed", "Character documents were imported explicitly.")
}

// RollbackCharacters restores the character documents recorded by the last
// explicit import. It refuses rather than overwriting a document changed since
// that import.
func (a *App) RollbackCharacters() error {
	store := a.characterStore()
	journal, ok := a.cutoverStore("character").ReadJournal()
	if !ok || len(journal.Items) == 0 {
		return fmt.Errorf("RollbackCharacters: no character import journal was found; nothing was rolled back")
	}
	entries := make([]characterstore.ImportResult, 0, len(journal.Items))
	for _, item := range journal.Items {
		current, err := store.Get(item.ID)
		if err != nil {
			return fmt.Errorf("RollbackCharacters: %s was changed or removed after import; nothing was rolled back", item.ID)
		}
		if item.Revision != "" && current.Revision != item.Revision {
			return fmt.Errorf("RollbackCharacters: %s changed after import; nothing was rolled back", item.ID)
		}
		entries = append(entries, characterstore.ImportResult{
			ID:       item.ID,
			Status:   item.Action,
			Revision: item.Revision,
			Previous: item.Previous,
		})
	}
	if err := store.Rollback(entries); err != nil {
		return fmt.Errorf("RollbackCharacters: %w", err)
	}
	if err := a.cutoverStore("character").MarkRolledBack("Character import was undone."); err != nil {
		return fmt.Errorf("RollbackCharacters: %w", err)
	}
	return nil
}

// recordWorldCutover persists the result of an explicit world/scenario import.
func (a *App) recordWorldCutover(domain string, imported int, items []domaincutover.Item) error {
	return a.cutoverStore(domain).Record(imported, items, "completed", "Legacy documents were imported explicitly.")
}

func (a *App) rollbackWorldCategory(domain string, store *worldstore.Store) (domaincutover.RollbackResult, error) {
	res := domaincutover.RollbackResult{Domain: domain, Status: "nothing_to_do", Items: []string{}}
	journal, ok := a.cutoverStore(domain).ReadJournal()
	if !ok || len(journal.Items) == 0 {
		res.Message = "No import journal was found; nothing was rolled back."
		return res, nil
	}
	for i := len(journal.Items) - 1; i >= 0; i-- {
		item := journal.Items[i]
		if item.ID == "" {
			res.Skipped++
			continue
		}
		if len(item.Hashes) == 0 {
			res.Failed++
			res.Items = append(res.Items, item.ID+": no journaled file hashes")
			continue
		}
		result, err := store.RemoveImported(item.ID, item.Revision, item.Hashes)
		if err != nil {
			res.Failed++
			res.Items = append(res.Items, item.ID+": "+err.Error())
			continue
		}
		switch result.Status {
		case "rolled_back":
			res.Removed++
		case "not_found":
			res.Skipped++
		default:
			res.Failed++
			message := result.Message
			if message == "" {
				message = result.Status
			}
			res.Items = append(res.Items, item.ID+": "+message)
		}
	}
	if res.Failed > 0 {
		res.Status = "rollback_partial"
		res.Message = "Some imported documents changed or could not be removed; the cutover marker was kept."
		return res, nil
	}
	if err := a.cutoverStore(domain).MarkRolledBack("Imported documents were removed; the browser copy is authoritative again."); err != nil {
		return res, err
	}
	res.Status = "rolled_back"
	res.Message = "Imported documents were removed; the browser copy is authoritative again."
	return res, nil
}

func (a *App) RollbackWorlds() (domaincutover.RollbackResult, error) {
	store, err := a.worldStore()
	if err != nil {
		return domaincutover.RollbackResult{}, err
	}
	return a.rollbackWorldCategory("world", store)
}

func (a *App) RollbackScenarios() (domaincutover.RollbackResult, error) {
	store, err := a.scenarioStore()
	if err != nil {
		return domaincutover.RollbackResult{}, err
	}
	return a.rollbackWorldCategory("scenario", store)
}

// RollbackImageGenerationImport restores the generation document from before
// the last explicit import, or clears it when the import created it.
func (a *App) RollbackImageGenerationImport() (domaincutover.RollbackResult, error) {
	res := domaincutover.RollbackResult{Domain: "image-generation", Status: "nothing_to_do", Items: []string{}}
	if a.images == nil || a.images.Settings() == nil {
		return res, fmt.Errorf("image settings store not initialised")
	}
	journal, ok := a.cutoverStore("image-generation").ReadJournal()
	if !ok || len(journal.Items) != 1 {
		res.Message = "No image generation import journal was found; nothing was rolled back."
		return res, nil
	}
	item := journal.Items[0]
	current, revision, err := a.images.Settings().LoadGeneration()
	if err != nil {
		return res, err
	}
	if item.Revision != "" && revision != item.Revision {
		res.Status = "rollback_conflict"
		res.Message = "Image generation settings changed after import; nothing was rolled back."
		return res, nil
	}
	if len(item.Previous) > 0 {
		restored, err := a.SaveImageGenerationSettings(item.Previous, revision)
		if err != nil {
			return res, err
		}
		res.Items = append(res.Items, fmt.Sprintf("restored revision %v", restored["revision"]))
	} else if len(current) > 0 {
		if _, err := a.images.Settings().ClearGeneration(); err != nil {
			return res, err
		}
		res.Items = append(res.Items, "cleared imported generation settings")
	}
	if err := a.cutoverStore("image-generation").MarkRolledBack("Image generation import was undone."); err != nil {
		return res, err
	}
	res.Status = "rolled_back"
	res.Message = "Image generation settings were restored to their pre-import state."
	return res, nil
}
