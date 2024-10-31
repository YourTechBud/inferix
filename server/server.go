package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
)

type Server struct {
	options Options

	// Workspace related stuff
	lock         sync.RWMutex
	workspaces   map[string]*Workspace
	configDriver ConfigDriver // This is to manage global configuration (like workspaces)
}

func New(opts Options) (*Server, error) {
	// Intialise a global configuration driver to manage workspaces
	var configDriver ConfigDriver
	if opts.ConfigDriver == ConfigDriverType_LibSQL {
		globalConfigDriverOptions := Options{ConfigDriver: opts.ConfigDriver, ConfigPath: filepath.Join(opts.ConfigPath, "inferix.db")}
		driver, err := initialiseConfigDriver(globalConfigDriverOptions)
		if err != nil {
			return nil, err
		}
		configDriver = driver
	}

	return &Server{
		options:      opts,
		workspaces:   make(map[string]*Workspace),
		configDriver: configDriver,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Create a default workspace for the default tenant
	if s.options.CreateDefaultWorkspace {
		if err := s.NewWorkspace(ctx, "default", "default"); err != nil {
			return err
		}
	}

	// Setup the router
	router := s.router()

	// Start the server
	fmt.Println("Starting server on port 4386")
	return http.ListenAndServe(":4386", router)
}
