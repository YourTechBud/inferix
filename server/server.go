package server

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm"
	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/go-chi/chi/v5"
)

func New() (chi.Router, error) {
	// Create all the modules
	llm, err := llm.New(llm.Config{
		Models: []models.Config{
			models.Config{Driver: "openai", TargetName: "hugging-quants/Meta-Llama-3.1-8B-Instruct-AWQ-INT4", Name: "Llama-3.1-8B-Instruct"},
			models.Config{Driver: "ollama", Name: "llama3.1"},
		},
		Backends: []backends.Config{
			backends.Config{Name: "ollama", BackendType: "ollama", Config: json.RawMessage(`{"baseUrl": "http://localhost:11434"}`)},
			backends.Config{Name: "openai", BackendType: "openai", Config: json.RawMessage(`{"baseUrl": "http://localhost:2242/v1"}`)},
		},
	})
	if err != nil {
		return nil, err
	}

	// Add the api routes
	router := chi.NewRouter()
	router.Mount("/api/v1", llm.Routes())

	return router, nil
}
