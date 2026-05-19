package rag

import (
	"math"
	"sort"
)

// Index provides in-memory cosine similarity search over chunk embeddings
type Index struct {
	chunks []Chunk
}

// NewIndex creates an empty index
func NewIndex() *Index {
	return &Index{}
}

// AddChunks adds chunks to the index
func (idx *Index) AddChunks(chunks []Chunk) {
	idx.chunks = append(idx.chunks, chunks...)
}

// Clear removes all chunks
func (idx *Index) Clear() {
	idx.chunks = nil
}

// Search finds the top-k most similar chunks to the query embedding
func (idx *Index) Search(queryEmbedding []float32, topK int) []RetrievedChunk {
	if len(idx.chunks) == 0 || len(queryEmbedding) == 0 {
		return nil
	}

	scores := make([]struct {
		chunk Chunk
		score float32
	}, 0, len(idx.chunks))

	for _, chunk := range idx.chunks {
		if len(chunk.Embedding) == 0 {
			continue
		}
		score := cosineSimilarity(queryEmbedding, chunk.Embedding)
		scores = append(scores, struct {
			chunk Chunk
			score float32
		}{chunk: chunk, score: score})
	}

	// Sort by score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Take topK
	if topK > len(scores) {
		topK = len(scores)
	}

	result := make([]RetrievedChunk, topK)
	for i := 0; i < topK; i++ {
		result[i] = RetrievedChunk{
			Chunk: scores[i].chunk,
			Score: scores[i].score,
		}
	}
	return result
}

// cosineSimilarity computes the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		ai := float64(a[i])
		bi := float64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}

// DeduplicateChunks removes near-duplicate chunks based on high similarity
func DeduplicateChunks(chunks []RetrievedChunk, threshold float32) []RetrievedChunk {
	if len(chunks) == 0 {
		return nil
	}

	var result []RetrievedChunk
	for _, c := range chunks {
		isDup := false
		for _, existing := range result {
			if c.DocID == existing.DocID && cosineSimilarity(c.Embedding, existing.Embedding) > threshold {
				isDup = true
				break
			}
		}
		if !isDup {
			result = append(result, c)
		}
	}
	return result
}
