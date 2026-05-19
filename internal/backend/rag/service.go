package rag

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/backend/types"
)

// Embedder defines the contract for text embedding generation
type Embedder interface {
	Embed(text string) ([]float32, error)
}

// Service orchestrates the RAG pipeline: ingest, index, retrieve
type Service struct {
	mu       sync.RWMutex
	storage  *Storage
	chunker  *Chunker
	index    *Index
	embedder Embedder

	projectRoot string
}

// NewService creates a new RAG service
func NewService(projectRoot string, embedder Embedder) *Service {
	return &Service{
		storage:     NewStorage(projectRoot),
		chunker:     NewChunker(),
		index:       NewIndex(),
		embedder:    embedder,
		projectRoot: projectRoot,
	}
}

// IngestDocument processes a raw text document: chunk, embed, store
func (s *Service) IngestDocument(doc *Document, rawText string) error {
	if doc.ID == "" {
		doc.ID = fmt.Sprintf("doc_%d", time.Now().UnixNano())
	}
	if doc.CreatedAt == "" {
		doc.CreatedAt = time.Now().Format(time.RFC3339)
	}
	if doc.Enabled {
		doc.Enabled = true
	}

	// 1. Chunk text
	chunks := s.chunker.ChunkText(doc.ID, rawText)
	if len(chunks) == 0 {
		return types.ErrInvalidInput("document contains no extractable text")
	}

	// 2. Embed chunks
	if s.embedder == nil {
		return types.ErrInternal("no embedder configured")
	}

	for i := range chunks {
		embedding, err := s.embedder.Embed(chunks[i].Text)
		if err != nil {
			return fmt.Errorf("failed to embed chunk %d: %w", i, err)
		}
		chunks[i].Embedding = embedding
	}

	doc.ChunkCount = len(chunks)

	// 3. Persist document metadata, raw text, and embeddings
	if err := s.storage.SaveDocument(doc, rawText); err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}
	if err := s.storage.SaveChunks(doc.ID, chunks); err != nil {
		return fmt.Errorf("failed to save chunks: %w", err)
	}

	// 4. Add to in-memory index
	s.mu.Lock()
	s.index.AddChunks(chunks)
	s.mu.Unlock()

	return nil
}

// LoadIndex rebuilds the in-memory index from persisted documents
func (s *Service) LoadIndex(scope DocumentScope, scopeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.index.Clear()

	docs, err := s.storage.ListDocuments(scope, scopeID)
	if err != nil {
		return fmt.Errorf("failed to list documents: %w", err)
	}

	for _, doc := range docs {
		if !doc.Enabled {
			continue
		}
		chunks, err := s.storage.LoadChunks(doc.ID)
		if err != nil {
			continue
		}
		s.index.AddChunks(chunks)
	}

	return nil
}

// Retrieve performs RAG retrieval for a query
func (s *Service) Retrieve(q Query) (*RetrievalResult, error) {
	if s.embedder == nil {
		return nil, types.ErrInternal("no embedder configured")
	}

	// Build query text from message + recent history
	queryText := buildQueryText(q)

	queryEmbedding, err := s.embedder.Embed(queryText)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Search index
	topK := q.MaxChunks
	if topK == 0 {
		topK = 5
	}
	chunks := s.index.Search(queryEmbedding, topK*2) // Search more to allow dedup filtering

	// Deduplicate
	chunks = DeduplicateChunks(chunks, 0.95)

	// Enforce max total tokens
	maxTokens := q.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1200
	}

	var result []RetrievedChunk
	totalTokens := 0
	for _, chunk := range chunks {
		chunkTokens := estimateTokens(chunk.Text)
		if totalTokens+chunkTokens > maxTokens {
			break
		}
		result = append(result, chunk)
		totalTokens += chunkTokens

		// Update last used for the document
		_ = s.storage.UpdateDocumentLastUsed(chunk.DocID)
	}

	return &RetrievalResult{
		Chunks:      result,
		TotalTokens: totalTokens,
	}, nil
}

// ListDocuments returns documents filtered by scope
func (s *Service) ListDocuments(scope DocumentScope, scopeID string) ([]Document, error) {
	return s.storage.ListDocuments(scope, scopeID)
}

// DeleteDocument removes a document from storage and index
func (s *Service) DeleteDocument(docID string) error {
	if err := s.storage.DeleteDocument(docID); err != nil {
		return err
	}

	// Rebuild index to remove deleted chunks
	s.mu.Lock()
	s.index.Clear()
	s.mu.Unlock()

	return nil
}

// EnableDocument toggles document enabled state
func (s *Service) EnableDocument(docID string, enabled bool) error {
	doc, err := s.storage.LoadDocument(docID)
	if err != nil {
		return err
	}
	doc.Enabled = enabled
	data, _ := json.MarshalIndent(doc, "", "  ")
	return os.WriteFile(s.storage.ragPath(docID+".json"), data, 0644)
}

// ExtractTextFromFile reads text from supported file types
func ExtractTextFromFile(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt", ".md", ".json", ".csv":
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
		return string(data), nil
	case ".pdf", ".epub", ".docx":
		return "", types.ErrInvalidInput("PDF/EPUB/DOCX extraction requires external tool; convert to TXT first")
	default:
		return "", types.ErrInvalidInput("unsupported file type: " + ext)
	}
}

// buildQueryText combines the current message and recent history into a query
func buildQueryText(q Query) string {
	var parts []string
	if len(q.MessageHistory) > 0 {
		// Take last 3 messages for context
		start := len(q.MessageHistory) - 3
		if start < 0 {
			start = 0
		}
		parts = append(parts, q.MessageHistory[start:]...)
	}
	parts = append(parts, q.Text)
	return strings.Join(parts, "\n")
}

// estimateTokens roughly estimates token count (1 token ≈ 4 chars for English)
func estimateTokens(text string) int {
	return int(math.Ceil(float64(len(text)) / 4.0))
}
