/**
 * Semantic Memory Clustering
 *
 * Groups similar memories into clusters using TF-IDF embeddings.
 * Enables memory organization, deduplication, and retrieval.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"fmt"
	"sort"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/types"
)

// Cluster represents a group of semantically similar memories
type Cluster struct {
	ID       string             `json:"id"`
	Label    string             `json:"label"`
	Memories []types.MemoryRecord `json:"memories"`
	Score    int                `json:"score"` // Average similarity score
}

// ClusterConfig controls clustering behavior
type ClusterConfig struct {
	SimilarityThreshold int  // Minimum similarity to join cluster (0-100)
	MinClusterSize      int  // Minimum memories to form a cluster
	MaxClusters         int  // Maximum clusters to return
}

// DefaultClusterConfig returns sensible defaults
func DefaultClusterConfig() ClusterConfig {
	return ClusterConfig{
		SimilarityThreshold: 65,
		MinClusterSize:      2,
		MaxClusters:         10,
	}
}

// ClusterMemories groups memories into semantic clusters
func ClusterMemories(memories []types.MemoryRecord, config ClusterConfig) []Cluster {
	if len(memories) < config.MinClusterSize {
		return nil
	}

	// Build embeddings for all memories
	type memoryEmbedding struct {
		memory    types.MemoryRecord
		embedding *lorebook.SimpleEmbedding
	}

	var items []memoryEmbedding
	for _, m := range memories {
		text := m.Title + " " + m.Summary + " " + m.Content
		items = append(items, memoryEmbedding{
			memory:    m,
			embedding: lorebook.NewEmbedding(text),
		})
	}

	// Greedy clustering: assign each memory to best-matching cluster
	var clusters []Cluster
	for _, item := range items {
		bestCluster := -1
		bestScore := 0

		for i, cluster := range clusters {
			score := clusterSimilarity(item.embedding, cluster)
			if score >= config.SimilarityThreshold && score > bestScore {
				bestCluster = i
				bestScore = score
			}
		}

		if bestCluster >= 0 {
			clusters[bestCluster].Memories = append(clusters[bestCluster].Memories, item.memory)
			clusters[bestCluster].Score = updateClusterScore(clusters[bestCluster])
		} else {
			clusters = append(clusters, Cluster{
				ID:       types.NewMemoryID(),
				Label:    generateClusterLabel(item.memory),
				Memories: []types.MemoryRecord{item.memory},
				Score:    100,
			})
		}
	}

	// Filter to clusters with minimum size
	var valid []Cluster
	for _, c := range clusters {
		if len(c.Memories) >= config.MinClusterSize {
			valid = append(valid, c)
		}
	}

	// Sort by score descending and limit
	sort.Slice(valid, func(i, j int) bool {
		return valid[i].Score > valid[j].Score
	})

	if len(valid) > config.MaxClusters {
		valid = valid[:config.MaxClusters]
	}

	return valid
}

// clusterSimilarity computes average similarity between an embedding and a cluster
func clusterSimilarity(emb *lorebook.SimpleEmbedding, cluster Cluster) int {
	if len(cluster.Memories) == 0 {
		return 0
	}

	total := 0
	for _, m := range cluster.Memories {
		text := m.Title + " " + m.Summary + " " + m.Content
		other := lorebook.NewEmbedding(text)
		total += emb.CosineSimilarity(other)
	}

	return total / len(cluster.Memories)
}

// updateClusterScore recalculates average intra-cluster similarity
func updateClusterScore(cluster Cluster) int {
	if len(cluster.Memories) < 2 {
		return 100
	}

	// Sample up to 10 pairs for efficiency
	var embeddings []*lorebook.SimpleEmbedding
	for _, m := range cluster.Memories {
		text := m.Title + " " + m.Summary + " " + m.Content
		embeddings = append(embeddings, lorebook.NewEmbedding(text))
	}

	total := 0
	count := 0
	for i := 0; i < len(embeddings) && i < 10; i++ {
		for j := i + 1; j < len(embeddings) && j < 10; j++ {
			total += embeddings[i].CosineSimilarity(embeddings[j])
			count++
		}
	}

	if count == 0 {
		return 100
	}
	return total / count
}

// generateClusterLabel creates a human-readable label from cluster contents
func generateClusterLabel(memory types.MemoryRecord) string {
	if memory.Title != "" {
		return fmt.Sprintf("About %s", memory.Title)
	}
	if len(memory.Tags) > 0 {
		return fmt.Sprintf("Related to %s", memory.Tags[0])
	}
	return "Memory Group"
}

// FindSimilarMemories returns memories similar to a query text
func FindSimilarMemories(memories []types.MemoryRecord, query string, threshold int) []types.MemoryRecord {
	queryEmb := lorebook.NewEmbedding(query)

	var similar []types.MemoryRecord
	for _, m := range memories {
		text := m.Title + " " + m.Summary + " " + m.Content
		memEmb := lorebook.NewEmbedding(text)
		if queryEmb.IsSemanticMatch(memEmb, threshold) {
			similar = append(similar, m)
		}
	}

	return similar
}
