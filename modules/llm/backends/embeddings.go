package backends

import (
	"context"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// CreateEmbeddings creates embeddings using the appropriate backend.
func (b *Backends) CreateEmbeddings(ctx context.Context, req types.EmbeddingRequest) (types.EmbeddingResponse, error) {
	// Get the model
	modelConfig, err := b.models.GetModel(req.Model)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Set the target name of the model
	req.Model = modelConfig.GetTarget()

	// Get the backend
	backend, err := b.getBackend(modelConfig.Backend)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Check if the backend supports embeddings
	if !backend.EnableEmbeddingsAPI() {
		return types.EmbeddingResponse{}, utils.NewStandardError(http.StatusBadRequest, "Embeddings API is disabled", "embeddings_disabled")
	}

	// Create embeddings
	response, err := backend.CreateEmbeddings(ctx, req)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Don't forget to set the model name back
	response.Model = modelConfig.GetID()

	return response, nil
}
