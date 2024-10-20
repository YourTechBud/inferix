package config

import (
	"github.com/YourTechBud/inferix/modules/llm/models"
)

// Config is a struct to configure the LLM module
type Config struct {
	Models   []models.Config `json:"models" yaml:"models"`
	Backends []BackendConfig `json:"backends" yaml:"backends"`
}
