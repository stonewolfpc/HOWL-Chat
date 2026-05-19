/**
 * Memory Deduplication Engine
 *
 * Detects and merges duplicate or near-duplicate memories.
 * Uses semantic similarity and title matching.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"strings"
	"time"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/types"
)

// DedupConfig controls deduplication behavior
type DedupConfig struct {
	SemanticThreshold int  // Minimum similarity to consider duplicate (0-100)
	TitleMatchBoost   int  // Extra points for title similarity
	PreferNewer       bool // When merging, keep newer or older
}

// DefaultDedupConfig returns sensible defaults
func DefaultDedupConfig() DedupConfig {
	return DedupConfig{
		SemanticThreshold: 75,
		TitleMatchBoost:   20,
		PreferNewer:       true,
	}
}

// DedupResult represents the outcome of deduplication
type DedupResult struct {
	Kept       []types.MemoryRecord // Surviving unique memories
	Merged     int                  // Number of merges performed
	Duplicates []string             // IDs of removed duplicates
}

// DeduplicateMemories scans a memory list and merges duplicates
func DeduplicateMemories(memories []types.MemoryRecord, config DedupConfig) DedupResult {
	if len(memories) < 2 {
		return DedupResult{Kept: memories}
	}

	var kept []types.MemoryRecord
	var removed []string
	merged := 0

	// For each memory, check if it's a duplicate of anything already kept
	for _, candidate := range memories {
		isDup := false
		for i, existing := range kept {
			if isDuplicate(candidate, existing, config) {
				kept[i] = mergeMemories(existing, candidate, config.PreferNewer)
				removed = append(removed, candidate.ID)
				merged++
				isDup = true
				break
			}
		}
		if !isDup {
			kept = append(kept, candidate)
		}
	}

	return DedupResult{
		Kept:       kept,
		Merged:     merged,
		Duplicates: removed,
	}
}

// isDuplicate checks if two memories are duplicates
func isDuplicate(a, b types.MemoryRecord, config DedupConfig) bool {
	// Exact title match = definite duplicate
	if strings.EqualFold(strings.TrimSpace(a.Title), strings.TrimSpace(b.Title)) &&
		a.Title != "" && a.Title != "Notable Event" {
		return true
	}

	// Semantic similarity check
	textA := a.Title + " " + a.Summary + " " + a.Content
	textB := b.Title + " " + b.Summary + " " + b.Content

	embA := lorebook.NewEmbedding(textA)
	embB := lorebook.NewEmbedding(textB)
	similarity := embA.CosineSimilarity(embB)

	// Title similarity boost
	if titlesSimilar(a.Title, b.Title) {
		similarity += config.TitleMatchBoost
	}

	return similarity >= config.SemanticThreshold
}

// titlesSimilar checks if two titles are similar
func titlesSimilar(a, b string) bool {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))
	if a == "" || b == "" {
		return false
	}
	if a == b {
		return true
	}
	// Check if one contains the other
	if strings.Contains(a, b) || strings.Contains(b, a) {
		return true
	}
	return false
}

// mergeMemories combines two memories, keeping the best of each
func mergeMemories(existing, incoming types.MemoryRecord, preferNewer bool) types.MemoryRecord {
	result := existing

	// Use the longer/more detailed content
	if len(incoming.Content) > len(existing.Content) {
		result.Content = incoming.Content
	}

	// Use the longer summary
	if len(incoming.Summary) > len(existing.Summary) {
		result.Summary = incoming.Summary
	}

	// Merge tags
	result.Tags = mergeTags(existing.Tags, incoming.Tags)

	// Keep highest confidence
	if incoming.Confidence > existing.Confidence {
		result.Confidence = incoming.Confidence
	}

	// Update timestamp
	now := time.Now().Format(time.RFC3339)
	result.UpdatedAt = now

	// Keep pinned status
	if incoming.Pinned {
		result.Pinned = true
	}

	// Increment trigger count (this memory has been "seen" again)
	result.TriggerCount++

	return result
}

// mergeTags combines two tag lists without duplicates
func mergeTags(a, b []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, t := range a {
		low := strings.ToLower(t)
		if !seen[low] {
			seen[low] = true
			result = append(result, t)
		}
	}

	for _, t := range b {
		low := strings.ToLower(t)
		if !seen[low] {
			seen[low] = true
			result = append(result, t)
		}
	}

	return result
}
