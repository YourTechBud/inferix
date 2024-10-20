package backends

import (
	"context"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// CreateEmbeddings creates embeddings using the appropriate backend.
func (b *Backends) CreateEmbeddings(ctx context.Context, req types.EmbeddingRequest) (types.EmbeddingResponse, error) {
	// Get the model
	modelConfig, err := b.models.GetModel(req.Model)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Set the target name of the model
	req.Model = modelConfig.GetTargetName()

	// Get the backend
	backend, err := b.getBackend(modelConfig.Driver)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Create embeddings
	response, err := backend.CreateEmbeddings(ctx, req)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Don't forget to set the model name back
	response.Model = modelConfig.GetName()

	return response, nil
}
