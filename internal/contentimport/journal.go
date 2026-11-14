package contentimport

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"howl-chat/internal/worldstore"
)

const markerVersion = 2

var (
	// ErrMarkerMissing reports a terminal import operation whose completion
	// marker is absent.
	ErrMarkerMissing = errors.New("completion marker is missing")
	// ErrMarkerMalformed reports a marker that cannot be decoded or has an
	// unsupported schema version.
	ErrMarkerMalformed = errors.New("completion marker is malformed")
	// ErrMarkerConflict reports a marker belonging to another operation,
	// category, plan, or terminal state.
	ErrMarkerConflict = errors.New("completion marker conflicts with the journal")
	// ErrMarkerStale reports a marker that survives for a non-completion state.
	ErrMarkerStale = errors.New("completion marker is stale")
)

func (s *Service) journalPath() string {
	return filepath.Join(s.configDir, string(s.category)+"-import-journal.json")
}

// markerPath is operation-scoped. A new operation therefore cannot overwrite
// another operation's completion marker merely by finishing.
func (s *Service) markerPath(operationID string) string {
	return filepath.Join(s.configDir, fmt.Sprintf("%s-import-marker-%s.json", s.category, operationID))
}

func (s *Service) loadJournal() journalFile {
	journal := journalFile{Version: 1, Operations: []Operation{}}
	data, err := os.ReadFile(s.journalPath())
	if err != nil {
		return journal
	}
	if err := json.Unmarshal(data, &journal); err != nil || journal.Operations == nil {
		return journalFile{Version: 1, Operations: []Operation{}}
	}
	return journal
}

func (s *Service) loadOperation(id string) (Operation, bool) {
	for _, operation := range s.loadJournal().Operations {
		if operation.OperationID == id {
			return operation, true
		}
	}
	return Operation{}, false
}

// Operations returns the durable journal operations for the page. It is a
// read-only copy and never triggers recovery.
func (s *Service) Operations() []Operation {
	journal := s.loadJournal()
	out := make([]Operation, len(journal.Operations))
	copy(out, journal.Operations)
	for i := range out {
		if err := s.verifyMarker(out[i]); err != nil {
			out[i].Status = OperationIntegrityFailure
			out[i].Message = err.Error()
		}
	}
	return out
}

// writeOperation appends or replaces one operation. The journal lock protects
// only this metadata file; it never authorizes a canonical write.
func (s *Service) writeOperation(operation Operation) error {
	s.journalMu.Lock()
	defer s.journalMu.Unlock()
	journal := s.loadJournal()
	found := false
	for i := range journal.Operations {
		if journal.Operations[i].OperationID == operation.OperationID {
			journal.Operations[i] = operation
			found = true
			break
		}
	}
	if !found {
		journal.Operations = append(journal.Operations, operation)
	}
	data, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return fmt.Errorf("encode import journal: %w", err)
	}
	if err := writeAtomicFile(s.journalPath(), data); err != nil {
		return fmt.Errorf("write import journal: %w", err)
	}
	return nil
}

// verifyMarker enforces agreement between the journal operation and its
// operation-scoped completion marker. A marker is required only after an
// import reaches completed or completed_partial. For every other state a
// surviving marker is stale and blocks recovery rather than being trusted.
func (s *Service) verifyMarker(operation Operation) error {
	if !worldstore.ValidID(operation.OperationID) {
		return fmt.Errorf("%w: invalid operation identity %q", ErrMarkerMalformed, operation.OperationID)
	}
	data, err := os.ReadFile(s.markerPath(operation.OperationID))
	if err != nil {
		if os.IsNotExist(err) {
			if operationNeedsMarker(operation.Status) {
				return fmt.Errorf("%w: operation %s is %s", ErrMarkerMissing, operation.OperationID, operation.Status)
			}
			return nil
		}
		return fmt.Errorf("%w: read: %v", ErrMarkerMalformed, err)
	}
	var m marker
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("%w: decode: %v", ErrMarkerMalformed, err)
	}
	if m.Version != markerVersion {
		return fmt.Errorf("%w: version %d, want %d", ErrMarkerMalformed, m.Version, markerVersion)
	}
	if m.Category != string(s.category) || m.OperationID != operation.OperationID {
		return fmt.Errorf("%w: marker belongs to category %q operation %q", ErrMarkerConflict, m.Category, m.OperationID)
	}
	if m.PlanHash != operation.PlanHash ||
		m.DestinationLibraryID != operation.DestinationLibraryID ||
		m.DestinationRoot != operation.DestinationRoot ||
		!sameStrings(m.ItemIDs, operationItemIDs(operation)) {
		return fmt.Errorf("%w: marker identity does not match operation %s", ErrMarkerStale, operation.OperationID)
	}
	if m.Status != operation.Status {
		return fmt.Errorf("%w: marker status %q, journal status %q", ErrMarkerConflict, m.Status, operation.Status)
	}
	if !operationNeedsMarker(operation.Status) {
		return fmt.Errorf("%w: marker survives for non-completion state %q", ErrMarkerStale, operation.Status)
	}
	return nil
}

// writeMarker creates a marker exclusively when absent. Replacing an existing
// marker is allowed only when it proves ownership by matching the operation,
// category, plan, destination, item set, and terminal status.
func (s *Service) writeMarker(operation Operation) error {
	if !operationNeedsMarker(operation.Status) {
		return nil
	}
	if !worldstore.ValidID(operation.OperationID) {
		return fmt.Errorf("%w: invalid operation identity %q", ErrMarkerMalformed, operation.OperationID)
	}

	s.markerMu.Lock()
	defer s.markerMu.Unlock()

	next := markerForOperation(operation)
	data, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return fmt.Errorf("encode import marker: %w", err)
	}

	existing, err := s.readMarker(operation.OperationID)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := writeNewFileExclusive(s.markerPath(operation.OperationID), data); err != nil {
			return fmt.Errorf("create import marker: %w", err)
		}
		return nil
	}
	if !sameMarker(existing, next) {
		return fmt.Errorf("%w: refusing to replace operation %s marker with %s", ErrMarkerConflict, existing.OperationID, operation.OperationID)
	}
	if err := writeAtomicFile(s.markerPath(operation.OperationID), data); err != nil {
		return fmt.Errorf("replace import marker: %w", err)
	}
	return nil
}

// syncMarker keeps terminal metadata aligned with the journal without ever
// inferring completion from a marker alone.
func (s *Service) syncMarker(operation Operation) error {
	if operationNeedsMarker(operation.Status) {
		return s.writeMarker(operation)
	}
	return s.removeMarker(operation.OperationID)
}

func (s *Service) readMarker(operationID string) (marker, error) {
	data, err := os.ReadFile(s.markerPath(operationID))
	if err != nil {
		return marker{}, err
	}
	var m marker
	if err := json.Unmarshal(data, &m); err != nil {
		return marker{}, fmt.Errorf("%w: decode: %v", ErrMarkerMalformed, err)
	}
	if m.Version != markerVersion {
		return marker{}, fmt.Errorf("%w: version %d, want %d", ErrMarkerMalformed, m.Version, markerVersion)
	}
	return m, nil
}

func (s *Service) removeMarker(operationID string) error {
	if !worldstore.ValidID(operationID) {
		return fmt.Errorf("%w: invalid operation identity %q", ErrMarkerMalformed, operationID)
	}
	data, err := os.ReadFile(s.markerPath(operationID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var m marker
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("%w: decode: %v", ErrMarkerMalformed, err)
	}
	if m.Version != markerVersion {
		return fmt.Errorf("%w: version %d, want %d", ErrMarkerMalformed, m.Version, markerVersion)
	}
	if m.Category != string(s.category) || m.OperationID != operationID {
		return fmt.Errorf("%w: refusing to remove operation %s marker", ErrMarkerConflict, m.OperationID)
	}
	if err := os.Remove(s.markerPath(operationID)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func operationNeedsMarker(status string) bool {
	switch status {
	case OperationCompleted, OperationCompletedPartial:
		return true
	}
	return false
}

func markerForOperation(operation Operation) marker {
	return marker{
		Version:              markerVersion,
		Category:             operation.Category,
		OperationID:          operation.OperationID,
		PlanHash:             operation.PlanHash,
		DestinationLibraryID: operation.DestinationLibraryID,
		DestinationRoot:      operation.DestinationRoot,
		ItemIDs:              operationItemIDs(operation),
		Status:               operation.Status,
		CompletedAt:          nowRFC3339(),
	}
}

func sameMarker(a, b marker) bool {
	return a.Version == b.Version &&
		a.Category == b.Category &&
		a.OperationID == b.OperationID &&
		a.PlanHash == b.PlanHash &&
		a.DestinationLibraryID == b.DestinationLibraryID &&
		a.DestinationRoot == b.DestinationRoot &&
		a.Status == b.Status &&
		sameStrings(a.ItemIDs, b.ItemIDs)
}

func operationItemIDs(operation Operation) []string {
	ids := make([]string, 0, len(operation.Items))
	for _, item := range operation.Items {
		if item.ItemID != "" {
			ids = append(ids, item.ItemID)
		}
	}
	sort.Strings(ids)
	return ids
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func writeNewFileExclusive(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	return f.Close()
}

func statPath(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

func writeAtomicFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
