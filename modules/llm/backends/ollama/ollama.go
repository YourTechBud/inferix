package ollama

import (
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/goccy/go-yaml"
)

// Ollama is a struct that handles all interactions with an Ollama backend.
type Ollama struct {
	BaseURL string `json:"base_url,omitempty"`

	// Other fields
	options config.BackendOptions `json:"-"`
}

// New creates a new Ollama backend with the provided configuration.
func New(config config.BackendConfig) (*Ollama, error) {
	// Parse the configuration.
	cfg := new(Ollama)
	if err := yaml.Unmarshal(config.Config, cfg); err != nil {
		return nil, err
	}

	// Don't forget to set the backend options.
	cfg.options = config.Options

	return cfg, nil
}

var _ types.Backend = (*Ollama)(nil)
