package lorebookstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"howl-chat/internal/library"
)

// ErrInvalidIdentityMap rejects malformed or conflicting character identities.
var ErrInvalidIdentityMap = errors.New("character identity map is invalid")

// OwnerIdentity is one authoritative character record supplied by the frontend
// library. Name is only a label; ID is the ownership identity.
type OwnerIdentity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OwnerMigrationItem reports one entry that needed attention.
type OwnerMigrationItem struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	OwnerID   string `json:"owner_id,omitempty"`
	OwnerName string `json:"owner_name,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

// OwnerMigrationResult reports every outcome separately. No count implies
// another; ambiguous and unknown entries are left untouched.
type OwnerMigrationResult struct {
	Status    string               `json:"status"`
	Version   string               `json:"version,omitempty"`
	Migrated  int                  `json:"migrated"`
	Already   int                  `json:"already"`
	Ambiguous int                  `json:"ambiguous"`
	Unknown   int                  `json:"unknown"`
	Failed    int                  `json:"failed"`
	Restored  int                  `json:"restored,omitempty"`
	Skipped   int                  `json:"skipped,omitempty"`
	Items     []OwnerMigrationItem `json:"items"`
	Message   string               `json:"message"`
}

type ownerMigrationMarker struct {
	Schema     int    `json:"schema"`
	Version    string `json:"version"`
	MigratedAt string `json:"migrated_at"`
	Migrated   int    `json:"migrated"`
	Already    int    `json:"already"`
	Ambiguous  int    `json:"ambiguous"`
	Unknown    int    `json:"unknown"`
}

type ownerMigrationJournalItem struct {
	ID                   string `json:"id"`
	LibraryID            string `json:"library_id"`
	Path                 string `json:"path"`
	OwnerID              string `json:"owner_id"`
	PreviousOwnerID      string `json:"previous_owner_id"`
	PreviousOwnerName    string `json:"previous_owner_name,omitempty"`
	HadPreviousOwnerName bool   `json:"had_previous_owner_name"`
}

type ownerMigrationJournal struct {
	Schema     int                         `json:"schema"`
	Version    string                      `json:"version"`
	MigratedAt string                      `json:"migrated_at"`
	Status     string                      `json:"status"`
	Items      []ownerMigrationJournalItem `json:"items"`
}

// resolveOwnerIdentities validates the map and groups names. Duplicate names are
// retained as ambiguous; duplicate IDs with different names are rejected.
func resolveOwnerIdentities(identities []OwnerIdentity) (map[string]string, map[string][]string, error) {
	ids := make(map[string]string, len(identities))
	byName := make(map[string][]string, len(identities))
	for _, identity := range identities {
		id := strings.TrimSpace(identity.ID)
		name := strings.TrimSpace(identity.Name)
		if id == "" || name == "" {
			return nil, nil, fmt.Errorf("%w: every identity needs both an id and a name", ErrInvalidIdentityMap)
		}
		if !validID.MatchString(id) {
			return nil, nil, fmt.Errorf("%w: invalid identity %q", ErrInvalidIdentityMap, id)
		}
		if existing, ok := ids[id]; ok {
			if existing != name {
				return nil, nil, fmt.Errorf("%w: identity %q is listed with two names", ErrInvalidIdentityMap, id)
			}
			continue
		}
		ids[id] = name
		byName[name] = append(byName[name], id)
	}
	if len(ids) == 0 {
		return nil, nil, fmt.Errorf("%w: no identities supplied", ErrInvalidIdentityMap)
	}
	for name := range byName {
		sort.Strings(byName[name])
	}
	return ids, byName, nil
}

// MigrateOwners rewrites legacy owner display names to stable character IDs.
// It only touches character-scope entries with a non-empty ownerId, never
// overwrites an ID that is already known, and never guesses an ambiguous or
// unknown name. Successful rewrites are journaled for rollback.
func (s *Store) MigrateOwners(identities []OwnerIdentity, version string) (OwnerMigrationResult, error) {
	res := OwnerMigrationResult{Status: "failed", Version: strings.TrimSpace(version), Items: []OwnerMigrationItem{}}
	if s.reg == nil {
		res.Message = "content library registry is not initialised"
		return res, nil
	}
	if res.Version == "" {
		return res, fmt.Errorf("%w: migration version is required", ErrInvalidIdentityMap)
	}
	ids, byName, err := resolveOwnerIdentities(identities)
	if err != nil {
		return res, err
	}
	if !s.Cutover() {
		return res, ErrNotCutover
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	items := make(map[string]Item)
	duplicates := make(map[string]bool)
	scanErrors := []string{}
	for _, view := range s.reg.List(library.Lorebooks) {
		if !view.Readable {
			continue
		}
		scanned, _, scanErr := s.scanRoot(view)
		if scanErr != nil {
			scanErrors = append(scanErrors, scanErr.Error())
			continue
		}
		for _, item := range scanned {
			if _, exists := items[item.ID]; exists {
				duplicates[item.ID] = true
				continue
			}
			items[item.ID] = item
		}
	}
	for id := range duplicates {
		delete(items, id)
		res.Failed++
		res.Items = append(res.Items, OwnerMigrationItem{
			ID:     id,
			Action: "failed",
			Detail: "the same identity exists in more than one library; left unchanged",
		})
	}
	// An unreadable root makes the migration incomplete even when every entry
	// that was visible succeeded. Do not advance the version marker in that case.
	res.Failed += len(scanErrors)

	ordered := make([]string, 0, len(items))
	for id := range items {
		ordered = append(ordered, id)
	}
	sort.Strings(ordered)

	now := time.Now().UTC().Format(time.RFC3339)
	journal, hadJournal := s.readOwnerJournal()
	if !hadJournal || journal.Status == "rolled_back" {
		journal = ownerMigrationJournal{Schema: 1, Version: res.Version, Items: []ownerMigrationJournalItem{}}
	}
	journal.Version = res.Version
	journal.MigratedAt = now
	journal.Status = "completed"
	journaled := make(map[string]bool, len(journal.Items))
	for _, item := range journal.Items {
		journaled[item.ID] = true
	}

	for _, id := range ordered {
		item := items[id]
		if item.Entry == nil {
			continue
		}
		scope := strings.ToLower(strings.TrimSpace(stringField(item.Entry, "scope")))
		if scope == "" {
			scope = strings.ToLower(strings.TrimSpace(stringField(item.Entry, "type")))
		}
		if scope != "character" {
			continue
		}
		owner := strings.TrimSpace(stringField(item.Entry, "ownerId"))
		if owner == "" {
			continue
		}
		if _, stable := ids[owner]; stable {
			res.Already++
			continue
		}
		candidates := byName[owner]
		switch {
		case len(candidates) > 1:
			res.Ambiguous++
			res.Items = append(res.Items, OwnerMigrationItem{
				ID:      id,
				Action:  "ambiguous",
				OwnerID: owner,
				Detail:  fmt.Sprintf("%q matches %d character records; left unchanged", owner, len(candidates)),
			})
			continue
		case len(candidates) == 0:
			res.Unknown++
			res.Items = append(res.Items, OwnerMigrationItem{
				ID:      id,
				Action:  "unknown",
				OwnerID: owner,
				Detail:  fmt.Sprintf("%q is not a known character; left unchanged", owner),
			})
			continue
		}
		if !item.Editable {
			res.Failed++
			res.Items = append(res.Items, OwnerMigrationItem{
				ID:      id,
				Action:  "failed",
				OwnerID: owner,
				Detail:  "the library holding this entry is read-only",
			})
			continue
		}

		stableID := candidates[0]
		updated := make(map[string]interface{}, len(item.Entry)+4)
		for key, value := range item.Entry {
			updated[key] = value
		}
		updated["ownerId"] = stableID
		updated["ownerName"] = owner
		updated["ownerMigratedFrom"] = owner
		updated["ownerMigrationVersion"] = res.Version

		payload, encodeErr := encodeDocument(id, updated)
		if encodeErr != nil {
			res.Failed++
			res.Items = append(res.Items, OwnerMigrationItem{ID: id, Action: "failed", Detail: encodeErr.Error()})
			continue
		}
		target, writeErr := s.reg.WriteFile(item.LibraryID, id+".json", payload)
		if writeErr != nil {
			res.Failed++
			res.Items = append(res.Items, OwnerMigrationItem{ID: id, Action: "failed", Detail: writeErr.Error()})
			continue
		}

		res.Migrated++
		res.Items = append(res.Items, OwnerMigrationItem{
			ID:        id,
			Action:    "migrated",
			OwnerID:   stableID,
			OwnerName: owner,
		})
		if !journaled[id] {
			journal.Items = append(journal.Items, ownerMigrationJournalItem{
				ID:                   id,
				LibraryID:            item.LibraryID,
				Path:                 target,
				OwnerID:              stableID,
				PreviousOwnerID:      owner,
				PreviousOwnerName:    strings.TrimSpace(stringField(item.Entry, "ownerName")),
				HadPreviousOwnerName: func() bool { _, ok := item.Entry["ownerName"]; return ok }(),
			})
			journaled[id] = true
		}
	}

	if res.Failed > 0 {
		journal.Status = "partial"
	}
	if len(journal.Items) > 0 {
		if err := s.writeOwnerJournal(journal); err != nil {
			return res, err
		}
	}
	if res.Failed == 0 {
		marker := ownerMigrationMarker{
			Schema:     1,
			Version:    res.Version,
			MigratedAt: now,
			Migrated:   res.Migrated,
			Already:    res.Already,
			Ambiguous:  res.Ambiguous,
			Unknown:    res.Unknown,
		}
		if err := s.writeOwnerMarker(marker); err != nil {
			return res, err
		}
	}

	switch {
	case res.Failed > 0:
		res.Status = "partial"
	case res.Migrated > 0 && (res.Ambiguous > 0 || res.Unknown > 0):
		res.Status = "completed_with_unresolved"
	case res.Migrated > 0:
		res.Status = "completed"
	case res.Ambiguous > 0 || res.Unknown > 0:
		res.Status = "needs_review"
	default:
		res.Status = "nothing_to_do"
	}
	if len(scanErrors) > 0 {
		res.Message = "Some libraries could not be read. " + strings.Join(scanErrors, " ")
	}
	res.Message += fmt.Sprintf(" Migrated %d, already stable %d, ambiguous %d, unknown %d, failed %d.",
		res.Migrated, res.Already, res.Ambiguous, res.Unknown, res.Failed)
	res.Message = strings.TrimSpace(res.Message)
	return res, nil
}

// RollbackOwners restores the ownership values recorded by MigrateOwners. It
// only restores an entry whose current owner is still the value the migration
// wrote, so later edits are never overwritten. Unrelated fields are preserved.
func (s *Store) RollbackOwners() (OwnerMigrationResult, error) {
	res := OwnerMigrationResult{Status: "failed", Items: []OwnerMigrationItem{}}
	if s.reg == nil {
		res.Message = "content library registry is not initialised"
		return res, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	journal, ok := s.readOwnerJournal()
	if !ok || len(journal.Items) == 0 {
		res.Status = "nothing_to_do"
		res.Message = "No owner migration journal was found; nothing was restored."
		return res, nil
	}
	if journal.Status == "rolled_back" {
		res.Status = "nothing_to_do"
		res.Message = "The recorded owner migration has already been rolled back."
		return res, nil
	}
	res.Version = journal.Version

	for _, record := range journal.Items {
		item, found := s.itemByID(record.ID)
		if !found {
			res.Skipped++
			res.Items = append(res.Items, OwnerMigrationItem{ID: record.ID, Action: "skipped", Detail: "entry no longer exists"})
			continue
		}
		if filepath.Clean(item.Path) != filepath.Clean(record.Path) || stringField(item.Entry, "ownerId") != record.OwnerID {
			res.Skipped++
			res.Items = append(res.Items, OwnerMigrationItem{ID: record.ID, Action: "skipped", Detail: "ownership changed after migration; left untouched"})
			continue
		}
		if !item.Editable {
			res.Failed++
			res.Items = append(res.Items, OwnerMigrationItem{ID: record.ID, Action: "failed", Detail: "the library holding this entry is read-only"})
			continue
		}

		updated := make(map[string]interface{}, len(item.Entry)+2)
		for key, value := range item.Entry {
			updated[key] = value
		}
		if record.PreviousOwnerID == "" {
			delete(updated, "ownerId")
		} else {
			updated["ownerId"] = record.PreviousOwnerID
		}
		if record.HadPreviousOwnerName {
			updated["ownerName"] = record.PreviousOwnerName
		} else {
			delete(updated, "ownerName")
		}
		delete(updated, "ownerMigratedFrom")
		delete(updated, "ownerMigrationVersion")

		payload, encodeErr := encodeDocument(record.ID, updated)
		if encodeErr != nil {
			res.Failed++
			res.Items = append(res.Items, OwnerMigrationItem{ID: record.ID, Action: "failed", Detail: encodeErr.Error()})
			continue
		}
		if _, writeErr := s.reg.WriteFile(item.LibraryID, record.ID+".json", payload); writeErr != nil {
			res.Failed++
			res.Items = append(res.Items, OwnerMigrationItem{ID: record.ID, Action: "failed", Detail: writeErr.Error()})
			continue
		}
		res.Restored++
		res.Items = append(res.Items, OwnerMigrationItem{ID: record.ID, Action: "restored", OwnerID: record.PreviousOwnerID})
	}

	if res.Failed > 0 {
		// Keep the journal retryable; a failed entry must not be declared rolled
		// back, and the version marker must stay until every write succeeds.
		journal.Status = "partial"
		if err := s.writeOwnerJournal(journal); err != nil {
			return res, err
		}
		res.Status = "partial"
		res.Message = fmt.Sprintf("Restored %d entries; %d could not be restored.", res.Restored, res.Failed)
		return res, nil
	}
	journal.Status = "rolled_back"
	if err := s.writeOwnerJournal(journal); err != nil {
		return res, err
	}
	if err := os.Remove(s.ownerMarkerPath); err != nil && !os.IsNotExist(err) {
		return res, err
	}
	res.Status = "rolled_back"
	res.Message = fmt.Sprintf("Restored %d ownership values; unrelated data was left untouched.", res.Restored)
	return res, nil
}

func (s *Store) readOwnerJournal() (ownerMigrationJournal, bool) {
	data, err := os.ReadFile(s.ownerJournalPath)
	if err != nil {
		return ownerMigrationJournal{}, false
	}
	var journal ownerMigrationJournal
	if err := json.Unmarshal(data, &journal); err != nil {
		return ownerMigrationJournal{}, false
	}
	return journal, true
}

func (s *Store) writeOwnerJournal(journal ownerMigrationJournal) error {
	data, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return fmt.Errorf("encode owner migration journal: %w", err)
	}
	return writeAtomic(s.ownerJournalPath, data)
}

func (s *Store) writeOwnerMarker(marker ownerMigrationMarker) error {
	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return fmt.Errorf("encode owner migration marker: %w", err)
	}
	return writeAtomic(s.ownerMarkerPath, data)
}
