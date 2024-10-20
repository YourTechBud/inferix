package tei

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/types"
)

// TEI is a struct that handles all interactions with the Text Embedding Inference backend.
type TEI struct {
	BaseURL string `json:"baseUrl"`
}

// New creates a new TEI struct with the provided configuration.
func New(config config.BackendConfig) (*TEI, error) {
	cfg := new(TEI)
	if err := json.Unmarshal(config.Config, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

var _ types.Backend = (*TEI)(nil)
