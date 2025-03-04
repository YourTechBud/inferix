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
	m.lock.RLock()
	defer m.lock.RUnlock()

	model, found := m.finalizedModels[modelName]
	if !found {
		return config.ModelConfig{}, utils.NewStandardError(http.StatusNotFound, fmt.Sprintf("model %s not found", modelName), "model_not_found")
	}

	return model, nil
}

// GetModels returns the list of models
func (m *Models) GetModels(ctx context.Context) ([]types.ModelObject, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	models := make([]types.ModelObject, 0, len(m.finalizedModels))

	for id, model := range m.finalizedModels {
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

// MergeModels merges the provided ModelObjects into the models map
func (m *Models) MergeModels(modelObjects []types.ModelObject) {
	// Create a new map for the merged models
	newModels := make(map[string]config.ModelConfig, len(modelObjects))

	// Convert ModelObjects to ModelConfig
	for _, obj := range modelObjects {
		// Create new ModelConfig from ModelObject
		modelConfig := config.ModelConfig{
			ID:             obj.ID,
			Backend:        obj.OwnedBy,
			Target:         obj.ID,
			Aliases:        []string{},
			DefaultOptions: config.DefaultModelOptions(),
		}

		// Apply any matching model settings
		for _, settings := range m.modelSettings {
			settings.ApplySettings(&modelConfig)
		}

		newModels[obj.ID] = modelConfig
	}

	// Finally, copy all configured models overriding similarly named models if any
	for k, v := range m.configuredModels {
		newModels[k] = v
	}

	// Update the mergedModels map
	m.lock.Lock()
	m.finalizedModels = newModels
	m.lock.Unlock()
}
