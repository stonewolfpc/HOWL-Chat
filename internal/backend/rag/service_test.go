package rag

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChunker(t *testing.T) {
	c := NewChunker()
	text := "This is the first sentence. Here is another one! And a third? Finally, the last sentence here."
	chunks := c.ChunkText("doc1", text)

	if len(chunks) == 0 {
		t.Fatal("Expected chunks, got none")
	}

	t.Logf("Produced %d chunks:", len(chunks))
	for i, chunk := range chunks {
		t.Logf("  Chunk %d: %q", i, chunk.Text)
	}
}

func TestCosineSimilarity(t *testing.T) {
	a := []float32{1, 0, 0}
	b := []float32{1, 0, 0}
	c := []float32{0, 1, 0}

	if cosineSimilarity(a, b) != 1.0 {
		t.Error("Expected similarity of identical vectors to be 1.0")
	}

	if cosineSimilarity(a, c) != 0.0 {
		t.Error("Expected similarity of orthogonal vectors to be 0.0")
	}
}

func TestStorage(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	doc := &Document{
		ID:         "test-doc-1",
		Name:       "Test Document",
		Scope:      ScopeWorld,
		ScopeID:    "Eternum",
		SourceType: SourceFile,
		SizeBytes:  100,
		Enabled:    true,
	}

	rawText := "This is a test document for RAG storage."

	if err := storage.SaveDocument(doc, rawText); err != nil {
		t.Fatalf("SaveDocument failed: %v", err)
	}

	loaded, err := storage.LoadDocument("test-doc-1")
	if err != nil {
		t.Fatalf("LoadDocument failed: %v", err)
	}
	if loaded.Name != "Test Document" {
		t.Errorf("Expected name 'Test Document', got %q", loaded.Name)
	}

	text, err := storage.LoadRawText("test-doc-1")
	if err != nil {
		t.Fatalf("LoadRawText failed: %v", err)
	}
	if text != rawText {
		t.Errorf("Expected raw text %q, got %q", rawText, text)
	}

	docs, err := storage.ListDocuments(ScopeWorld, "Eternum")
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}
	if len(docs) != 1 {
		t.Errorf("Expected 1 document, got %d", len(docs))
	}

	_ = storage.DeleteDocument("test-doc-1")
	_, err = storage.LoadDocument("test-doc-1")
	if err == nil {
		t.Error("Expected document to be deleted")
	}
}

func TestKeywordEmbedder(t *testing.T) {
	emb := NewKeywordEmbedder()

	vec1, err := emb.Embed("hello world test")
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	if len(vec1) == 0 {
		t.Error("Expected non-empty vector")
	}

	vec2, err := emb.Embed("hello world example")
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	sim := cosineSimilarity(vec1, vec2)
	if sim <= 0 {
		t.Errorf("Expected positive similarity for related texts, got %f", sim)
	}
}

func TestExtractTextFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	txtPath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(txtPath, []byte("Hello from file."), 0644)

	text, err := ExtractTextFromFile(txtPath)
	if err != nil {
		t.Fatalf("ExtractTextFromFile failed: %v", err)
	}
	if text != "Hello from file." {
		t.Errorf("Expected 'Hello from file.', got %q", text)
	}

	_, err = ExtractTextFromFile("test.pdf")
	if err == nil {
		t.Error("Expected error for unsupported PDF")
	}
}
