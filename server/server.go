package server

import (
	"github.com/YourTechBud/inferix/modules/llm"
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/utils"
	"github.com/go-chi/chi/v5"
)

func New(configFilePath string) (chi.Router, error) {
	// Read the yaml file from the path provided
	var cfg config.Config
	if err := utils.ReadYAMLFile(configFilePath, &cfg); err != nil {
		return nil, err
	}

	// Create all the modules
	llm, err := llm.New(cfg)
	if err != nil {
		return nil, err
	}

	// Add the api routes
	router := chi.NewRouter()
	router.Mount("/api/v1", llm.Routes())

	return router, nil
}
