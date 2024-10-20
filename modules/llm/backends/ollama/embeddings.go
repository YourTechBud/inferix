package ollama

import (
	"context"
	"fmt"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// CreateEmbeddings creates embeddings using the Ollama backend.
func (backend *Ollama) CreateEmbeddings(ctx context.Context, req types.EmbeddingRequest) (types.EmbeddingResponse, error) {
	// Check if the embeddings API is enabled
	if !backend.options.EnableEmbeddingsAPI {
		return types.EmbeddingResponse{}, utils.NewStandardError(http.StatusBadRequest, "Embeddings API is disabled", "embeddings_disabled")
	}

	// Create ollama request
	ollamaReq := OllamaEmbeddingsRequest{
		EmbeddingRequest: req,
		Truncate:         false,
	}

	// Make the request
	ollamaResponse, err := utils.MakeHTTPRequest[OllamaEmbeddingsResponse, utils.StandardError](ctx, http.MethodPost, fmt.Sprintf("%s/api/embed", backend.BaseURL), ollamaReq)
	if err != nil {
		err = utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_call_error")
		return types.EmbeddingResponse{}, err
	}

	// Check if we have an error
	if ollamaResponse.HasError() {
		// Set the status of the error to the status of the response
		ollamaResponse.Error.Status = ollamaResponse.Status

		// Return the error
		return types.EmbeddingResponse{}, ollamaResponse.Error
	}

	// Convert the response
	embeddings := make([]types.Embedding, len(ollamaResponse.Data.Embeddings))
	for i, embedding := range ollamaResponse.Data.Embeddings {
		embeddings[i] = types.Embedding{
			Index:     uint32(i),
			Embedding: embedding,
			Object:    types.EmbeddingObject_Embedding,
		}
	}

	return types.EmbeddingResponse{
		Model:  ollamaResponse.Data.Model,
		Data:   embeddings,
		Usage:  types.EmbeddingUsage{PromptTokens: ollamaResponse.Data.PromptEvalCount, TotalTokens: ollamaResponse.Data.PromptEvalCount},
		Object: types.EmbeddingResponseObject_List,
	}, nil
}
