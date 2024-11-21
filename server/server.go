package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/YourTechBud/inferix/utils"
)

type Server struct {
	options       Options
	defaultConfig map[string]any

	// Workspace related stuff
	lock       sync.RWMutex
	workspaces map[string]*Workspace
}

func New(opts Options) (*Server, error) {
	return &Server{
		options:    opts,
		workspaces: make(map[string]*Workspace),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Create a default workspace for the default tenant
	if s.options.CreateDefaultWorkspace {
		if err := s.NewWorkspace(ctx, "default", "default"); err != nil {
			return err
		}
	}

	// Load the default configuration if provided
	if s.options.DefaultConfigPath != "" {
		if err := utils.ReadYAMLFile(s.options.DefaultConfigPath, &s.defaultConfig); err != nil {
			return err
		}
	}

	// Setup the router
	router := s.router()

	// Start the server
	fmt.Println("Starting server on port 4386")
	return http.ListenAndServe(":4386", router)
}
