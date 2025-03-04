package llm

import (
	"encoding/json"
	"sync"

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

	// Goroutine management
	done chan struct{}
	wg   sync.WaitGroup
}

// New creates a new LLM struct
func New(workspaceCtx *utils.WorkspaceContext, cfg json.RawMessage) (utils.Module, error) {
	// Unmarshal the configuration
	config := new(config.Config)
	if err := json.Unmarshal(cfg, config); err != nil {
		return nil, err
	}

	// Create a new models struct
	models := models.New(config.ModelConfigs, config.ModelSettings, workspaceCtx.Storage)

	// Create a new backends struct
	backends, err := backends.New(config.BackendConfigs, models)
	if err != nil {
		return nil, err
	}

	llm := &LLM{
		models:   models,
		backends: backends,
		routes:   initializeRoutes(models, backends),
		done:     make(chan struct{}),
	}

	// Start model polling
	llm.startModelPolling()

	// Return the module
	return llm, nil
}

// Close closes the module
func (llm *LLM) Close() error {
	// Signal all goroutines to stop
	close(llm.done)

	// Wait for the model polling goroutine to finish
	llm.wg.Wait()

	// Close the router
	llm.routes = nil

	return nil
}
