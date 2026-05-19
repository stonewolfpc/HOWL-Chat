package rag

import (
	"strings"
	"unicode"
)

// Chunker splits text into overlapping chunks
type Chunker struct {
	ChunkSize    int
	ChunkOverlap int
}

// NewChunker creates a default chunker with 512 tokens size and 128 overlap
func NewChunker() *Chunker {
	return &Chunker{
		ChunkSize:    512,
		ChunkOverlap: 128,
	}
}

// ChunkText splits raw text into chunks with overlap
func (c *Chunker) ChunkText(docID string, text string) []Chunk {
	// Normalize text
	text = normalizeText(text)
	if len(text) == 0 {
		return nil
	}

	// Split into sentences first for semantic boundaries
	sentences := splitSentences(text)
	if len(sentences) == 0 {
		return nil
	}

	var chunks []Chunk
	var current strings.Builder
	currentLen := 0
	idx := 0

	for _, sentence := range sentences {
		sentenceLen := len(sentence)
		if sentenceLen > c.ChunkSize {
			// Single sentence is longer than chunk size — split by words
			if currentLen > 0 {
				chunks = append(chunks, makeChunk(docID, idx, current.String()))
				idx++
				current.Reset()
				currentLen = 0
			}
			wordChunks := c.splitByWords(docID, sentence, &idx)
			chunks = append(chunks, wordChunks...)
			continue
		}

		if currentLen+sentenceLen > c.ChunkSize && currentLen > 0 {
			// Finish current chunk
			chunks = append(chunks, makeChunk(docID, idx, current.String()))
			idx++

			// Start new chunk with overlap from previous chunk
			overlap := getOverlap(current.String(), c.ChunkOverlap)
			current.Reset()
			current.WriteString(overlap)
			currentLen = len(overlap)
		}

		if currentLen > 0 {
			current.WriteString(" ")
			currentLen++
		}
		current.WriteString(sentence)
		currentLen += sentenceLen
	}

	if currentLen > 0 {
		chunks = append(chunks, makeChunk(docID, idx, current.String()))
	}

	return chunks
}

// splitByWords splits a long text into word-aligned chunks
func (c *Chunker) splitByWords(docID string, text string, idx *int) []Chunk {
	words := strings.Fields(text)
	var chunks []Chunk
	var current strings.Builder
	currentLen := 0

	for _, word := range words {
		wordLen := len(word) + 1 // +1 for space
		if currentLen+wordLen > c.ChunkSize && currentLen > 0 {
			chunks = append(chunks, makeChunk(docID, *idx, strings.TrimSpace(current.String())))
			(*idx)++
			current.Reset()
			currentLen = 0
		}
		if currentLen > 0 {
			current.WriteString(" ")
			currentLen++
		}
		current.WriteString(word)
		currentLen += len(word)
	}

	if currentLen > 0 {
		chunks = append(chunks, makeChunk(docID, *idx, strings.TrimSpace(current.String())))
		(*idx)++
	}
	return chunks
}

func makeChunk(docID string, index int, text string) Chunk {
	return Chunk{
		ID:    docID + "_chunk_" + string(rune('0'+index)),
		DocID: docID,
		Text:  text,
		Index: index,
	}
}

// normalizeText cleans text for chunking
func normalizeText(text string) string {
	// Replace multiple whitespace with single space
	var result strings.Builder
	var lastRune rune
	for _, r := range text {
		if unicode.IsSpace(r) {
			if lastRune != ' ' {
				result.WriteRune(' ')
				lastRune = ' '
			}
			continue
		}
		result.WriteRune(r)
		lastRune = r
	}
	return strings.TrimSpace(result.String())
}

// splitSentences splits text into sentences
func splitSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	for i, r := range text {
		current.WriteRune(r)
		if isSentenceEnd(r) {
			// Look ahead to see if next non-space char is uppercase
			if i+1 < len(text) {
				next := text[i+1:]
				next = strings.TrimLeftFunc(next, unicode.IsSpace)
				if len(next) > 0 && (unicode.IsUpper(rune(next[0])) || unicode.IsDigit(rune(next[0]))) {
					sentences = append(sentences, strings.TrimSpace(current.String()))
					current.Reset()
				}
			}
		}
	}

	if current.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(current.String()))
	}
	return sentences
}

func isSentenceEnd(r rune) bool {
	return r == '.' || r == '!' || r == '?'
}

// getOverlap returns the last n characters (or nearest sentence boundary) for overlap
func getOverlap(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	// Try to find a sentence boundary within the overlap window
	start := len(text) - maxLen
	for i := start; i < len(text); i++ {
		if isSentenceEnd(rune(text[i])) {
			return strings.TrimSpace(text[i+1:])
		}
	}
	// Fallback: just take the last maxLen chars from a word boundary
	for i := start; i < len(text); i++ {
		if text[i] == ' ' {
			return strings.TrimSpace(text[i+1:])
		}
	}
	return text[start:]
}
