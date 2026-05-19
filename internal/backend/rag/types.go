package rag

// DocumentScope defines the visibility scope of a document
type DocumentScope string

const (
	ScopeWorld    DocumentScope = "world"
	ScopeScenario DocumentScope = "scenario"
	ScopeChat     DocumentScope = "chat"
)

// SourceType defines how the document was ingested
type SourceType string

const (
	SourceFile SourceType = "file"
	SourceURL  SourceType = "url"
	SourceNote SourceType = "note"
)

// Document represents an uploaded knowledge document
type Document struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Scope      DocumentScope `json:"scope"`
	ScopeID    string        `json:"scopeId"` // world name, scenario name, or chat id
	SourceType SourceType    `json:"sourceType"`
	SourcePath string        `json:"sourcePath"` // original file path or URL
	SizeBytes  int64         `json:"sizeBytes"`
	ChunkCount int           `json:"chunkCount"`
	LastUsed   string        `json:"lastUsed"`
	CreatedAt  string        `json:"createdAt"`
	Enabled    bool          `json:"enabled"`
}

// Chunk represents a text segment with its embedding
type Chunk struct {
	ID        string    `json:"id"`
	DocID     string    `json:"docId"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"` // position in document
}

// RetrievedChunk represents a chunk returned from similarity search
type RetrievedChunk struct {
	Chunk
	Score float32 `json:"score"`
}

// RetrievalResult holds the outcome of a RAG query
type RetrievalResult struct {
	Chunks      []RetrievedChunk `json:"chunks"`
	TotalTokens int              `json:"totalTokens"`
}

// Query represents a RAG search query
type Query struct {
	Text           string
	WorldName      string
	ScenarioName   string
	ChatID         string
	MaxTokens      int
	MaxChunks      int
	MessageHistory []string
}
