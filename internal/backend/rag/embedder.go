package rag

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// LlamaEmbedder generates embeddings via llama-server HTTP API
type LlamaEmbedder struct {
	client interface{ Embedding(string) ([]float32, error) }
}

// NewLlamaEmbedder creates an embedder backed by llama-server
func NewLlamaEmbedder(client interface{ Embedding(string) ([]float32, error) }) *LlamaEmbedder {
	return &LlamaEmbedder{client: client}
}

// Embed generates an embedding vector for the given text
func (e *LlamaEmbedder) Embed(text string) ([]float32, error) {
	if e.client == nil {
		return nil, fmt.Errorf("llama client not available")
	}
	// Normalize: truncate if too long (llama embedding has token limit)
	normalized := normalizeForEmbedding(text)
	return e.client.Embedding(normalized)
}

// KeywordEmbedder is a fallback embedder using keyword frequency vectors
// Used when no LLM embedding endpoint is available
type KeywordEmbedder struct {
	vocab map[string]int
}

// NewKeywordEmbedder creates a simple keyword-based embedder
func NewKeywordEmbedder() *KeywordEmbedder {
	return &KeywordEmbedder{vocab: make(map[string]int)}
}

// Embed creates a sparse keyword frequency vector
func (e *KeywordEmbedder) Embed(text string) ([]float32, error) {
	words := tokenize(text)
	if len(e.vocab) == 0 {
		// First call: build vocab from this text
		for _, w := range words {
			e.vocab[w] = len(e.vocab)
		}
	}

	vec := make([]float32, len(e.vocab))
	for _, w := range words {
		if idx, ok := e.vocab[w]; ok {
			vec[idx]++
		}
	}

	// Normalize to unit vector
	norm := float32(0)
	for _, v := range vec {
		norm += v * v
	}
	if norm > 0 {
		norm = float32(math.Sqrt(float64(norm)))
		for i := range vec {
			vec[i] /= norm
		}
	}

	return vec, nil
}

func normalizeForEmbedding(text string) string {
	// Simple normalization: truncate to ~2000 chars to avoid token limits
	if len(text) > 2000 {
		text = text[:2000]
	}
	return strings.TrimSpace(text)
}

func tokenize(text string) []string {
	var words []string
	var current strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(unicode.ToLower(r))
		} else if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}
