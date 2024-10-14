package ollama

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// Ollama is a struct that handles all interactions with an Ollama backend.
type Ollama struct {
	BaseURL string `json:"baseUrl,omitempty"`
}

// New creates a new Ollama backend with the provided configuration.
func New(config json.RawMessage) (*Ollama, error) {
	// Parse the configuration.
	cfg := new(Ollama)
	if err := json.Unmarshal(config, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

var _ types.Backend = (*Ollama)(nil)
