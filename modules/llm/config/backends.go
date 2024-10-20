package config

import (
	"github.com/YourTechBud/inferix/utils"
)

// BackendConfig is a struct for backend configuration
type BackendConfig struct {
	Name        string           `json:"name"`
	BackendType string           `json:"type"`
	Config      utils.RawMessage `json:"config"`
	Options     BackendOptions   `json:"options"`
}

// BackendOptions is a struct for backend options
type BackendOptions struct {
	InjectFnCallPrompt  bool `json:"inject_fn_call_prompt"`
	EnableEmbeddingsAPI bool `json:"enable_embeddings_api"`
}
