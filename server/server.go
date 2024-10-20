package server

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm"
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/go-chi/chi/v5"
)

func New() (chi.Router, error) {
	// Create all the modules
	llm, err := llm.New(config.Config{
		Models: []models.Config{
			{
				Driver:     "openai",
				TargetName: "hugging-quants/Meta-Llama-3.1-8B-Instruct-AWQ-INT4",
				Name:       "Llama-3.1-8B-Instruct",
			},
			{
				Driver: "ollama",
				Name:   "llama3.1",
			},
			{
				Driver:     "tei",
				Name:       "bge-large-en-v1.5",
				Aliases:    []string{"default"},
				TargetName: "BAAI/bge-large-en-v1.5",
			},
		},
		Backends: []config.BackendConfig{
			{
				Name:        "ollama",
				BackendType: "ollama",
				Config:      json.RawMessage(`{"baseUrl": "http://localhost:11434"}`),
			},
			{
				Name:        "openai",
				BackendType: "openai",
				Options: config.BackendOptions{
					InjectFnCallPrompt: true,
				},
				Config: json.RawMessage(`{"baseUrl": "http://localhost:2242/v1"}`),
			},
			{
				Name:        "tei",
				BackendType: "tei",
				Config:      json.RawMessage(`{"baseUrl": "http://localhost:8080/v1"}`),
			},
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
