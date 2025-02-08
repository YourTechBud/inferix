package ollama

import (
	"context"
	"fmt"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// GetModels retrieves the list of available models from the Ollama API
func (o *Ollama) GetModels(ctx context.Context) ([]types.ModelObject, error) {
	// Make HTTP request to Ollama API
	resp, err := utils.MakeHTTPRequest[OllamaModelsResponse, any](ctx, "GET", fmt.Sprintf("%s/api/tags", o.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models from Ollama: %w", err)
	}

	if resp.HasError() {
		return nil, fmt.Errorf("failed to fetch models from Ollama: %v", resp.Error)
	}

	// Convert Ollama models to ModelObject format
	models := make([]types.ModelObject, len(resp.Data.Models))
	for i, model := range resp.Data.Models {
		models[i] = types.ModelObject{
			ID:      model.Name,
			Created: model.ModifiedAt.Unix(),
			Object:  "model",
			OwnedBy: "ollama",
		}
	}

	return models, nil
}
