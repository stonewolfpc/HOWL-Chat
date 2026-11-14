/**
 * Candidate Queue Manager
 *
 * Persists detected memory events awaiting user review.
 * Supports add, approve, reject, and query operations.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"howl-chat/internal/backend/types"
)

// CandidateStorage manages the candidate queue file
type CandidateStorage struct {
	projectRoot string
	chatID      string
}

// NewCandidateStorage creates a storage accessor for a chat session
func NewCandidateStorage(chatID string) *CandidateStorage {
	root, _ := os.Getwd()
	return &CandidateStorage{
		projectRoot: root,
		chatID:      chatID,
	}
}

// loadCandidates reads the candidate queue from disk
func (cs *CandidateStorage) loadCandidates() ([]types.MemoryCandidate, error) {
	path := cs.candidatesPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.MemoryCandidate{}, nil
		}
		return nil, fmt.Errorf("failed to read candidates: %w", err)
	}

	var candidates []types.MemoryCandidate
	if err := json.Unmarshal(data, &candidates); err != nil {
		return nil, fmt.Errorf("failed to parse candidates: %w", err)
	}
	return candidates, nil
}

// saveCandidates writes the candidate queue to disk
func (cs *CandidateStorage) saveCandidates(candidates []types.MemoryCandidate) error {
	chatDir := filepath.Join(cs.projectRoot, "Chats", cs.chatID)
	path := filepath.Join(chatDir, "candidates.json")

	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}

	data, err := json.MarshalIndent(candidates, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal candidates: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to save candidates: %w", err)
	}
	return nil
}

// candidatesPath returns the full file path
func (cs *CandidateStorage) candidatesPath() string {
	return filepath.Join(cs.projectRoot, "Chats", cs.chatID, "candidates.json")
}

// --- Service-level API ---

// AddCandidates appends new candidates to the queue, skipping duplicates
func (s *Service) AddCandidates(candidates []types.MemoryCandidate) error {
	storage := NewCandidateStorage(s.chatID)
	existing, err := storage.loadCandidates()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	for _, c := range existing {
		seen[c.SourceMessage] = true
	}

	added := 0
	for _, c := range candidates {
		if !seen[c.SourceMessage] {
			existing = append(existing, c)
			seen[c.SourceMessage] = true
			added++
		}
	}

	if added > 0 {
		return storage.saveCandidates(existing)
	}
	return nil
}

// GetCandidates returns candidates filtered by status
func (s *Service) GetCandidates(status string) ([]types.MemoryCandidate, error) {
	storage := NewCandidateStorage(s.chatID)
	all, err := storage.loadCandidates()
	if err != nil {
		return nil, err
	}

	if status == "" || status == "all" {
		return all, nil
	}

	var filtered []types.MemoryCandidate
	for _, c := range all {
		if c.Status == status {
			filtered = append(filtered, c)
		}
	}
	return filtered, nil
}

// ApproveCandidate converts a pending candidate into an active memory
func (s *Service) ApproveCandidate(candidateID string) (*types.MemoryRecord, error) {
	storage := NewCandidateStorage(s.chatID)
	candidates, err := storage.loadCandidates()
	if err != nil {
		return nil, err
	}

	var target *types.MemoryCandidate
	for i := range candidates {
		if candidates[i].ID == candidateID {
			target = &candidates[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("candidate not found: %s", candidateID)
	}

	// Convert to MemoryRecord
	now := time.Now().Format(time.RFC3339)
	record := &types.MemoryRecord{
		ID:         target.ID,
		Title:      target.Title,
		Summary:    target.Summary,
		Content:    target.Content,
		Scope:      target.Scope,
		OwnerID:    target.OwnerID,
		OwnerKind:  target.OwnerKind,
		Tags:       target.Tags,
		Pinned:     false,
		Status:     "active",
		Confidence: target.Confidence,
		Source:     "auto",
		CreatedAt:  target.CreatedAt,
		UpdatedAt:  now,
	}

	// Save to memories
	memStorage := NewMemoryStorage(s.chatID)
	memories, err := memStorage.loadMemories()
	if err != nil {
		return nil, err
	}
	memories = append(memories, *record)
	if err := memStorage.saveMemories(memories); err != nil {
		return nil, err
	}

	// Mark candidate as approved
	target.Status = "approved"
	target.ReviewedAt = now
	if err := storage.saveCandidates(candidates); err != nil {
		return nil, err
	}

	return record, nil
}

// RejectCandidate marks a candidate as rejected
func (s *Service) RejectCandidate(candidateID string) error {
	storage := NewCandidateStorage(s.chatID)
	candidates, err := storage.loadCandidates()
	if err != nil {
		return err
	}

	found := false
	for i := range candidates {
		if candidates[i].ID == candidateID {
			candidates[i].Status = "rejected"
			candidates[i].ReviewedAt = time.Now().Format(time.RFC3339)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("candidate not found: %s", candidateID)
	}

	return storage.saveCandidates(candidates)
}

// GetPendingCandidateCount returns number of pending candidates for UI badge
func (s *Service) GetPendingCandidateCount() (int, error) {
	pending, err := s.GetCandidates("pending")
	if err != nil {
		return 0, err
	}
	return len(pending), nil
}
