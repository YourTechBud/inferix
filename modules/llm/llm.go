package llm

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/models"
)

// LLM is a module to all kinds of interaction with Large Language Models
type LLM struct {
	// Internal stuff
	Models   *models.Models
	Backends *backends.Backends
}

// New creates a new LLM struct
func New(cfg json.RawMessage) (*LLM, error) {
	// Unmarshal the configuration
	config := new(config.Config)
	if err := json.Unmarshal(cfg, config); err != nil {
		return nil, err
	}

	// Create a new models struct
	models := models.New(config.Models)

	// Create a new backends struct
	backends, err := backends.New(config.Backends, models)
	if err != nil {
		return nil, err
	}

	// Return the module
	return &LLM{
		Backends: backends,
	}, nil
}
