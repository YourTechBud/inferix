package llm

import (
	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/models"
)

// Config is a struct to configure the LLM module
type Config struct {
	Models   []models.Config   `json:"models"`
	Backends []backends.Config `json:"backends"`
}
