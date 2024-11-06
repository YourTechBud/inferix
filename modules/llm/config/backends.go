package config

import (
	"encoding/json"
)

type (
	// BackendConfig is a struct for backend configuration
	BackendConfig struct {
		ID          string          `json:"id" validate:"required"`
		BackendType string          `json:"type" validate:"required"`
		Config      json.RawMessage `json:"config" validate:"required"`
		Options     BackendOptions  `json:"options"`
	}

	// BackendOptions is a struct for backend options
	BackendOptions struct {
		InjectFnCallPrompt  bool `json:"inject_fn_call_prompt"`
		EnableEmbeddingsAPI bool `json:"enable_embeddings_api"`
	}
)

func (c *BackendConfig) GetID() string {
	return c.ID
}
