package apikeys

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/utils"
)

type Module struct {
	middlewares []utils.HTTPMiddleware
}

func New(cfg json.RawMessage) (utils.Module, error) {
	// Unmarshal the configuration
	config := new(Config)
	if err := json.Unmarshal(cfg, config); err != nil {
		return nil, err
	}

	middlewares := []utils.HTTPMiddleware{
		initializeMiddleware(config),
	}
	// Return the module
	return &Module{
		middlewares: middlewares,
	}, nil
}

func (module *Module) Close() error {
	// Close the middleware
	module.middlewares = nil

	// Return nil
	return nil
}
