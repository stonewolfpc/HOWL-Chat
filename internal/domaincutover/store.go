// Package domaincutover persists the explicit ownership transition for domains
// that predate the content-library authority. A marker means the domain's
// canonical store owns the data; a journal records exactly what the last
// explicit import changed so the transition can be undone without guessing.
package domaincutover

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const version = 1

// Item records one document changed by an import.
type Item struct {
	ID       string                 `json:"id"`
	Action   string                 `json:"action"`
	Revision string                 `json:"revision,omitempty"`
	Path     string                 `json:"path,omitempty"`
	Hashes   map[string]string      `json:"hashes,omitempty"`
	Previous map[string]interface{} `json:"previous,omitempty"`
}

// Journal is the durable record of the last explicit import for one domain.
type Journal struct {
	Version    int    `json:"version"`
	Domain     string `json:"domain"`
	ImportedAt string `json:"imported_at"`
	Status     string `json:"status"`
	Items      []Item `json:"items"`
	Message    string `json:"message,omitempty"`
}

// Status is the read-only state exposed to the renderer.
type Status struct {
	Domain            string `json:"domain"`
	Cutover           bool   `json:"cutover"`
	CutoverAt         string `json:"cutover_at,omitempty"`
	OperationStatus   string `json:"operation_status,omitempty"`
	Imported          int    `json:"imported"`
	RollbackAvailable bool   `json:"rollback_available"`
	Message           string `json:"message,omitempty"`
}

// RollbackResult reports exactly what an undo changed.
type RollbackResult struct {
	Domain  string   `json:"domain"`
	Status  string   `json:"status"`
	Removed int      `json:"removed"`
	Skipped int      `json:"skipped"`
	Failed  int      `json:"failed"`
	Items   []string `json:"items,omitempty"`
	Message string   `json:"message,omitempty"`
}

// Store owns the marker and journal files for one domain.
type Store struct {
	mu          sync.Mutex
	domain      string
	markerPath  string
	journalPath string
}

// New builds the marker store for a domain name such as "character".
func New(configDir, domain string) *Store {
	name := strings.TrimSpace(strings.ToLower(domain))
	return &Store{
		domain:      name,
		markerPath:  filepath.Join(configDir, name+"-cutover.json"),
		journalPath: filepath.Join(configDir, name+"-import-journal.json"),
	}
}

// ReadJournal returns the last import record.
func (s *Store) ReadJournal() (Journal, bool) {
	data, err := os.ReadFile(s.journalPath)
	if err != nil {
		return Journal{}, false
	}
	var journal Journal
	if err := json.Unmarshal(data, &journal); err != nil || journal.Domain != s.domain {
		return Journal{}, false
	}
	return journal, true
}

// Status reports whether ownership has been explicitly transferred.
func (s *Store) Status() Status {
	status := Status{Domain: s.domain}
	data, err := os.ReadFile(s.markerPath)
	if err == nil {
		var marker struct {
			Version   int    `json:"version"`
			Domain    string `json:"domain"`
			CutoverAt string `json:"cutover_at"`
			Imported  int    `json:"imported"`
			Operation string `json:"operation_status"`
		}
		if json.Unmarshal(data, &marker) == nil && marker.Domain == s.domain {
			status.Cutover = true
			status.CutoverAt = marker.CutoverAt
			status.Imported = marker.Imported
			status.OperationStatus = marker.Operation
		}
	}
	if journal, ok := s.ReadJournal(); ok {
		status.OperationStatus = journal.Status
		status.RollbackAvailable = status.Cutover && len(journal.Items) > 0 && journal.Status != "rolled_back"
		status.Message = journal.Message
	}
	return status
}

// Record writes the journal first, then the marker. A crash in between leaves
// the journal visible but ownership unclaimed; the next explicit import can
// safely replace it.
func (s *Store) Record(imported int, items []Item, operationStatus, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	journal := Journal{
		Version:    version,
		Domain:     s.domain,
		ImportedAt: time.Now().UTC().Format(time.RFC3339),
		Status:     operationStatus,
		Items:      items,
		Message:    message,
	}
	if journal.Status == "" {
		journal.Status = "completed"
	}
	if err := writeJSON(s.journalPath, journal); err != nil {
		return fmt.Errorf("write %s import journal: %w", s.domain, err)
	}
	marker := map[string]interface{}{
		"version":          version,
		"domain":           s.domain,
		"cutover_at":       journal.ImportedAt,
		"imported":         imported,
		"operation_status": journal.Status,
	}
	if err := writeJSON(s.markerPath, marker); err != nil {
		return fmt.Errorf("write %s cutover marker: %w", s.domain, err)
	}
	return nil
}

// MarkRolledBack retires the marker and keeps the journal as history.
func (s *Store) MarkRolledBack(message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	journal, ok := s.ReadJournal()
	if !ok {
		return errors.New("no import journal was found")
	}
	journal.Status = "rolled_back"
	journal.Message = message
	if err := writeJSON(s.journalPath, journal); err != nil {
		return fmt.Errorf("write %s import journal: %w", s.domain, err)
	}
	if err := os.Remove(s.markerPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s cutover marker: %w", s.domain, err)
	}
	return nil
}

func writeJSON(path string, value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
