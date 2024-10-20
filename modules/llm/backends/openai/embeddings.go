package openai

import (
	"context"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func (backend *OpenAI) CreateEmbeddings(ctx context.Context, req types.EmbeddingRequest) (types.EmbeddingResponse, error) {
	// Check if the embeddings API is enabled
	if !backend.options.EnableEmbeddingsAPI {
		return types.EmbeddingResponse{}, utils.NewStandardError(http.StatusBadRequest, "Embeddings API is disabled", "embeddings_disabled")
	}

	// Make the request
	openaiResponse, err := utils.MakeHTTPRequest[types.EmbeddingResponse, utils.StandardError](ctx, http.MethodPost, backend.BaseURL+"/embeddings", req)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Check if we have an error
	if openaiResponse.HasError() {
		// Set the status of the error to the status of the response
		openaiResponse.Error.Status = openaiResponse.Status

		// Return the error
		return types.EmbeddingResponse{}, openaiResponse.Error
	}

	return *openaiResponse.Data, nil
}
