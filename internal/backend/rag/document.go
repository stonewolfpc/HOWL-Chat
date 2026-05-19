package rag

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/backend/types"
)

const ragDir = "RAG"

// Storage handles file-based RAG document persistence
type Storage struct {
	projectRoot string
}

// NewStorage creates a new RAG storage instance
func NewStorage(projectRoot string) *Storage {
	return &Storage{projectRoot: projectRoot}
}

func (s *Storage) ragPath(parts ...string) string {
	return filepath.Join(s.projectRoot, ragDir, filepath.Join(parts...))
}

func (s *Storage) ensureRAGDir() error {
	return os.MkdirAll(s.ragPath(), 0755)
}

// SaveDocument persists a document metadata file and its raw text
func (s *Storage) SaveDocument(doc *Document, rawText string) error {
	if err := s.ensureRAGDir(); err != nil {
		return fmt.Errorf("failed to create RAG dir: %w", err)
	}

	// Write metadata
	metaPath := s.ragPath(doc.ID + ".json")
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}
	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write document metadata: %w", err)
	}

	// Write raw text
	textPath := s.ragPath(doc.ID + ".txt")
	if err := os.WriteFile(textPath, []byte(rawText), 0644); err != nil {
		return fmt.Errorf("failed to write raw text: %w", err)
	}

	return nil
}

// LoadDocument loads a document by ID
func (s *Storage) LoadDocument(docID string) (*Document, error) {
	metaPath := s.ragPath(docID + ".json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, types.WrapError(types.ErrorCodeIO, "document not found: "+docID, err)
		}
		return nil, fmt.Errorf("failed to read document: %w", err)
	}

	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}
	return &doc, nil
}

// LoadRawText loads the raw text content of a document
func (s *Storage) LoadRawText(docID string) (string, error) {
	textPath := s.ragPath(docID + ".txt")
	data, err := os.ReadFile(textPath)
	if err != nil {
		return "", fmt.Errorf("failed to read raw text: %w", err)
	}
	return string(data), nil
}

// ListDocuments returns all documents, optionally filtered by scope
func (s *Storage) ListDocuments(scope DocumentScope, scopeID string) ([]Document, error) {
	if err := s.ensureRAGDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.ragPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read RAG dir: %w", err)
	}

	var docs []Document
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(s.ragPath(entry.Name()))
		if err != nil {
			continue
		}
		var doc Document
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}
		if scope != "" && doc.Scope != scope {
			continue
		}
		if scopeID != "" && doc.ScopeID != scopeID {
			continue
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// DeleteDocument removes a document and its associated files
func (s *Storage) DeleteDocument(docID string) error {
	metaPath := s.ragPath(docID + ".json")
	textPath := s.ragPath(docID + ".txt")
	embPath := s.ragPath(docID + ".emb")

	_ = os.Remove(metaPath)
	_ = os.Remove(textPath)
	_ = os.Remove(embPath)
	return nil
}

// SaveChunks persists chunk embeddings for a document
func (s *Storage) SaveChunks(docID string, chunks []Chunk) error {
	embPath := s.ragPath(docID + ".emb")
	data, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chunks: %w", err)
	}
	return os.WriteFile(embPath, data, 0644)
}

// LoadChunks loads chunk embeddings for a document
func (s *Storage) LoadChunks(docID string) ([]Chunk, error) {
	embPath := s.ragPath(docID + ".emb")
	data, err := os.ReadFile(embPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read chunks: %w", err)
	}
	var chunks []Chunk
	if err := json.Unmarshal(data, &chunks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chunks: %w", err)
	}
	return chunks, nil
}

// UpdateDocumentLastUsed updates the LastUsed timestamp
func (s *Storage) UpdateDocumentLastUsed(docID string) error {
	doc, err := s.LoadDocument(docID)
	if err != nil {
		return err
	}
	doc.LastUsed = time.Now().Format(time.RFC3339)
	data, _ := json.MarshalIndent(doc, "", "  ")
	return os.WriteFile(s.ragPath(docID+".json"), data, 0644)
}
