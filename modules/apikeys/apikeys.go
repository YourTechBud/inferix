package apikeys

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/utils"
)

type Module struct {
	apiKeys          map[string]APIKey
	workspaceStorage *utils.WorkspaceStorage
}

func New(workspaceCtx *utils.WorkspaceContext, cfg json.RawMessage) (utils.Module, error) {
	// Unmarshal the configuration
	config := new(Config)
	if err := json.Unmarshal(cfg, config); err != nil {
		return nil, err
	}

	// Store the API keys in a map for easy access
	apiKeys := make(map[string]APIKey, len(config.APIKeys))
	for _, apiKey := range config.APIKeys {
		apiKeys[apiKey.ID] = apiKey
	}

	// Return the module
	return &Module{
		apiKeys:          apiKeys,
		workspaceStorage: workspaceCtx.Storage,
	}, nil
}

func (module *Module) Close() error {
	// Return nil
	return nil
}
