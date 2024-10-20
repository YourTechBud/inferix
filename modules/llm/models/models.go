package models

import (
	"strings"
)

// Models stores the model configurations
type Models struct {
	models map[string]Config
}

// New initializes the models map with the given configurations
func New(modelConfigs []Config) *Models {
	models := make(map[string]Config)
	for _, model := range modelConfigs {
		// Trim whitespace
		model.Name = strings.TrimSpace(model.Name)

		if model.DefaultOptions == nil {
			model.DefaultOptions = DefaultModelOptions()
		}
		models[model.Name] = model

		// Add aliases
		for _, alias := range model.Aliases {
			alias = strings.TrimSpace(alias)
			models[alias] = model
		}
	}
	return &Models{models}
}
