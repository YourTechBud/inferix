package tei

import "github.com/YourTechBud/inferix/modules/llm/types"

// EmbeddingRequest is the request struct for the Embed API
type EmbeddingRequest struct {
	Inputs   types.EmbeddingInput `json:"inputs"`
	Truncate bool                 `json:"truncate"`
}

// EmbeddingResponse is the response struct for the Embed API
type EmbeddingResponse [][]float32
