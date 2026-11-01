package main

import (
	"fmt"

	"howl-chat/internal/backend/rag"
)

// UploadRAGDocumentText ingests a document whose text the frontend already holds.
// The browser file input only exposes content, not a path, so this avoids the
// dead end of handing the backend a filename it cannot open.
func (a *App) UploadRAGDocumentText(name string, scope string, scopeID string, sourceType string, content string) (string, error) {
	if a.ragService == nil {
		return "", fmt.Errorf("RAG service not initialized")
	}
	if len(content) == 0 {
		return "", fmt.Errorf("document %q has no readable text", name)
	}
	doc := &rag.Document{
		Name:       name,
		Scope:      rag.DocumentScope(scope),
		ScopeID:    scopeID,
		SourceType: rag.SourceType(sourceType),
		SizeBytes:  int64(len(content)),
		Enabled:    true,
	}
	if err := a.ragService.IngestDocument(doc, content); err != nil {
		return "", fmt.Errorf("failed to ingest document: %w", err)
	}
	return doc.ID, nil
}
