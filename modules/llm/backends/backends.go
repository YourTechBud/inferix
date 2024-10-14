package backends

import (
	"errors"

	"github.com/YourTechBud/inferix/modules/llm/backends/ollama"
	"github.com/YourTechBud/inferix/modules/llm/backends/openai"
	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/YourTechBud/inferix/modules/llm/types"
)

// Backends is a wrapper for all the backends
type Backends struct {
	models   *models.Models
	backends map[string]types.Backend
}

// New creates a new Backends struct
func New(backends []Config, models *models.Models) (*Backends, error) {
	backendsMap := make(map[string]types.Backend, len(backends))
	for _, backendConfig := range backends {
		switch backendConfig.BackendType {
		case "openai":
			backend, err := openai.New(backendConfig.Config)
			if err != nil {
				return nil, err
			}
			backendsMap[backendConfig.Name] = backend

		case "ollama":
			backend, err := ollama.New(backendConfig.Config)
			if err != nil {
				return nil, err
			}
			backendsMap[backendConfig.Name] = backend

		default:
			return nil, errors.New("unknown backend type")
		}
	}

	return &Backends{
		models:   models,
		backends: backendsMap,
	}, nil
}
