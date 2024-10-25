package models

import (
	"strings"

	"github.com/YourTechBud/inferix/modules/llm/config"
)

// Models stores the model configurations
type Models struct {
	models map[string]config.ModelConfig
}

// New initializes the models map with the given configurations
func New(modelConfigs []config.ModelConfig) *Models {
	models := make(map[string]config.ModelConfig)
	for _, model := range modelConfigs {
		// Trim whitespace
		model.ID = strings.TrimSpace(model.ID)

		if model.DefaultOptions == nil {
			model.DefaultOptions = config.DefaultModelOptions()
		}
		models[model.ID] = model

		// Add aliases
		for _, alias := range model.Aliases {
			alias = strings.TrimSpace(alias)
			models[alias] = model
		}
	}
	return &Models{models}
}
