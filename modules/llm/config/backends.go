package config

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

type (
	// BackendConfig is a struct for backend configuration
	BackendConfig struct {
		ID          string               `json:"id" validate:"required"`
		BackendType string               `json:"type" validate:"required"`
		Config      json.RawMessage      `json:"config" validate:"required"`
		Options     types.BackendOptions `json:"options"`
	}
)

func (c *BackendConfig) GetID() string {
	return c.ID
}
