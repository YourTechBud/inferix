package models

import (
	"strings"
	"sync"

	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/utils"
)

// Models stores the model configurations
type Models struct {
	lock sync.RWMutex

	configuredModels map[string]config.ModelConfig
	finalizedModels  map[string]config.ModelConfig
	modelSettings    []config.ModelSettings
	storage          *utils.WorkspaceStorage
}

// New initializes the models map with the given configurations
func New(modelConfigs []config.ModelConfig, modelSettings []config.ModelSettings, storage *utils.WorkspaceStorage) *Models {
	models := make(map[string]config.ModelConfig)
	for _, model := range modelConfigs {
		// Trim whitespace
		model.ID = strings.TrimSpace(model.ID)

		if model.DefaultOptions == nil {
			model.DefaultOptions = config.DefaultModelOptions()
		}

		// Apply any matching model settings
		for _, settings := range modelSettings {
			settings.ApplySettings(&model)
		}

		models[model.ID] = model

		// Add aliases
		for _, alias := range model.Aliases {
			alias = strings.TrimSpace(alias)
			models[alias] = model
		}
	}
	return &Models{
		configuredModels: models,
		finalizedModels:  models, // Initially, the finalized models are the same as the configured models
		modelSettings:    modelSettings,
		storage:          storage,
	}
}
