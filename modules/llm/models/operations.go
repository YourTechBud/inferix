package models

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// GetModel returns the model configuration for the given model name
func (m *Models) GetModel(modelName string) (config.ModelConfig, error) {

	model, found := m.models[modelName]
	if !found {
		return config.ModelConfig{}, utils.NewStandardError(http.StatusNotFound, fmt.Sprintf("model %s not found", modelName), "model_not_found")
	}

	return model, nil
}

// GetModels returns the list of models
func (m *Models) GetModels(ctx context.Context) ([]types.ModelObject, error) {
	models := make([]types.ModelObject, 0, len(m.models))

	for id, model := range m.models {
		model := types.ModelObject{
			ID:      id,
			Created: utils.ProcessStartTime,
			Object:  "model",
			OwnedBy: model.Backend,
		}

		// Try reading the model from the storage just to get the created time. Lol.
		// But we don't care cause SQLite is fast.
		resource, found, err := m.storage.GetResource(ctx, "models", id)
		if err != nil {
			log.Default().Println("Error getting model", err)
			return nil, utils.NewStandardError(http.StatusInternalServerError, "Error getting model", "error_getting_model")
		}
		if found {
			model.Created = resource.CreatedAt
		}

		models = append(models, model)
	}

	return models, nil
}
