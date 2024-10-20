package types

// EmbeddingRequest represents a request to create an embedding.
type EmbeddingRequest struct {
	Input EmbeddingInput `json:"input"`
	Model string         `json:"model"`
}

// EmbeddingResponse represents the response to an embedding request.
type EmbeddingResponse struct {
	Model  string                  `json:"model"`
	Data   []Embedding             `json:"data"`
	Usage  EmbeddingUsage          `json:"usage"`
	Object EmbeddingResponseObject `json:"object"`
}

// EmbeddingResponseObject represents the type of object in the embedding response.
type EmbeddingResponseObject string

const (
	EmbeddingResponseObject_List EmbeddingResponseObject = "list"
)

// Embedding represents the embedding data in the response.
type Embedding struct {
	Index     uint32          `json:"index"`
	Embedding []float32       `json:"embedding"`
	Object    EmbeddingObject `json:"object"`
}

// EmbeddingObject represents the type of object in the embedding.
type EmbeddingObject string

const (
	EmbeddingObject_Embedding EmbeddingObject = "embedding"
)

// EmbeddingUsage represents usage information for the embedding request.
type EmbeddingUsage struct {
	PromptTokens uint64 `json:"prompt_tokens"`
	TotalTokens  uint64 `json:"total_tokens"`
}

// EmbeddingInput is the input for the embedding API. This can either be a string or a list of strings.
type EmbeddingInput any
