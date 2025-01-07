package backends

import (
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/YourTechBud/inferix/modules/llm/types"
)

// Backends is a wrapper for all the backends
type Backends struct {
	models   *models.Models
	backends map[string]types.Backend
}

// New creates a new Backends struct
func New(backends []config.BackendConfig, models *models.Models) (*Backends, error) {
	backendsMap := make(map[string]types.Backend, len(backends))
	for _, backendConfig := range backends {
		backend, err := createBackend(backendConfig)
		if err != nil {
			return nil, err
		}
		backendsMap[backendConfig.ID] = backend
	}

	return &Backends{
		models:   models,
		backends: backendsMap,
	}, nil
}
