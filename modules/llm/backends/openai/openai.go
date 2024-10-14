package openai

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// OpenAI is a struct that handles all interactions with an OpenAI API compatible backend.
type OpenAI struct {
	APIKey  string `json:"apiKey,omitempty"`
	BaseURL string `json:"baseUrl,omitempty"`
}

// New creates a new OpenAI struct with the given API key and base URL.
func New(config json.RawMessage) (*OpenAI, error) {
	// Parse the configuration.
	cfg := new(OpenAI)
	if err := json.Unmarshal(config, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

var _ types.Backend = (*OpenAI)(nil)
