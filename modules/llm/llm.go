package llm

import (
	"encoding/json"

	"github.com/go-chi/chi/v5"

	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/YourTechBud/inferix/utils"
)

// LLM is a module to all kinds of interaction with Large Language Models
type LLM struct {
	// Internal stuff
	models   *models.Models
	backends *backends.Backends
	routes   chi.Router
}

// New creates a new LLM struct
func New(workspaceCtx *utils.WorkspaceContext, cfg json.RawMessage) (utils.Module, error) {
	// Unmarshal the configuration
	config := new(config.Config)
	if err := json.Unmarshal(cfg, config); err != nil {
		return nil, err
	}

	// Create a new models struct
	models := models.New(config.Models, workspaceCtx.Storage)

	// Create a new backends struct
	backends, err := backends.New(config.Backends, models)
	if err != nil {
		return nil, err
	}

	// Return the module
	return &LLM{
		models:   models,
		backends: backends,
		routes:   initializeRoutes(models, backends),
	}, nil
}

// Close closes the module
func (llm *LLM) Close() error {
	// TODO: close all the backends

	// Close the router
	llm.routes = nil

	// Return nil
	return nil
}
