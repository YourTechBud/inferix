package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"

	llmconfig "github.com/YourTechBud/inferix/modules/llm/config"
)

type Server struct {
	configDriver ConfigDriver

	chanUpdateConfig chan struct{}

	moduleLock sync.RWMutex
	modules    map[string]chi.Router
}

func New(opts Options) (*Server, error) {
	// Initialise the config driver
	configDriver, err := initialiseConfigDriver(opts)
	if err != nil {
		return nil, err
	}

	return &Server{
		configDriver:     configDriver,
		chanUpdateConfig: make(chan struct{}, 5),
		modules:          make(map[string]chi.Router),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Configure the modules
	if err := s.configureModules(ctx); err != nil {
		return err
	}

	// Start the module config updater
	go s.configUpdater()

	// Setup the router
	router := chi.NewRouter()
	router.Mount("/inferix/v1/llm", s.createModuleRouter("llm"))
	router.Mount("/inferix/v1/config/llm", s.createConfigRoutes("llm", llmconfig.GetConfigurationResources()))

	// Setup the global config route
	router.Get("/inferix/v1/config", s.getGlobalConfigHandler())

	// Start the server
	fmt.Println("Starting server on port 4386")
	return http.ListenAndServe(":4386", router)
}
