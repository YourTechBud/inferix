package openai

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/types"
)

// OpenAI is a struct that handles all interactions with an OpenAI API compatible backend.
type OpenAI struct {
	APIKey  string `json:"apiKey,omitempty"`
	BaseURL string `json:"baseUrl,omitempty"`

	// Other fields
	options config.BackendOptions `json:"-"`
}

// New creates a new OpenAI struct with the given API key and base URL.
func New(config config.BackendConfig) (*OpenAI, error) {
	// Parse the configuration.
	cfg := new(OpenAI)
	if err := json.Unmarshal(config.Config, cfg); err != nil {
		return nil, err
	}

	// Don't forget to set the backend options.
	cfg.options = config.Options

	return cfg, nil
}

var _ types.Backend = (*OpenAI)(nil)
