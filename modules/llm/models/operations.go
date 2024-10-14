package models

import (
	"fmt"
	"net/http"

	"github.com/YourTechBud/inferix/utils"
)

// GetModel returns the model configuration for the given model name
func (m *Models) GetModel(modelName string) (Config, error) {

	model, found := m.models[modelName]
	if !found {
		return Config{}, utils.NewStandardError(http.StatusNotFound, fmt.Sprintf("model %s not found", modelName), "model_not_found")
	}

	return model, nil
}
