package server

import (
	"context"
	"fmt"
	"path"
	"sync"

	"github.com/go-chi/chi/v5"
)

// Workspace is a map of modules
type Workspace struct {
	lock         sync.RWMutex
	modules      map[string]Module
	configDriver ConfigDriver

	// Channel to signal config update
	chanUpdateConfig chan struct{}

	// Router for the workspace
	router chi.Router
}

// NewWorkspace creates a new workspace
func (s *Server) NewWorkspace(ctx context.Context, tenant, workspace string) error {
	// Acquire the server lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the workspace key
	workspaceKey := getWorkspaceKey(tenant, workspace)

	// Check if the workspace already exists
	if _, ok := s.workspaces[workspaceKey]; ok {
		return fmt.Errorf("workspace already exists: %s", workspace)
	}

	// Setup the path for the workspace
	configPath := s.options.ConfigPath
	if s.options.ConfigDriver == ConfigDriverType_LibSQL {
		// For libsql driver, the config path is grouped by tenant and workspace
		configPath = path.Join(s.options.ConfigPath, tenant, workspace, "config.db")
	}

	// Create a config driver for the workspace
	configDriver, err := initialiseConfigDriver(Options{
		ConfigDriver:      s.options.ConfigDriver,
		ConfigPath:        configPath,
		DefaultConfigPath: s.options.DefaultConfigPath,
	})
	if err != nil {
		return err
	}

	// Create and return the workspace
	s.workspaces[workspaceKey] = &Workspace{
		// Create the modules map
		modules: make(map[string]Module),

		// Set the config driver
		configDriver: configDriver,

		// Create the channels
		chanUpdateConfig: make(chan struct{}, 5),
	}

	// Start the workspace
	if err := s.workspaces[workspaceKey].Start(ctx); err != nil {
		return err
	}

	return nil
}

// RemoveWorkspace removes a workspace
func (s *Server) RemoveWorkspace(tenant, workspace string) error {
	// Acquire the server lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Get the workspace key
	workspaceKey := getWorkspaceKey(tenant, workspace)

	// Check if the workspace exists
	w, p := s.workspaces[workspaceKey]
	if !p {
		return nil
	}

	// Stop the workspace
	// TODO: Deal with the error in a better way
	_ = w.Stop()

	// Remove the workspace
	delete(s.workspaces, workspaceKey)

	return nil
}

// Start begins the workspace operations
func (workspace *Workspace) Start(ctx context.Context) error {
	// Configure the modules
	if err := workspace.configureModules(ctx); err != nil {
		return err
	}

	// Start the module config updater
	go workspace.configUpdater()

	// Start the router
	workspace.intializeRouter()

	return nil
}

// Stop stops the workspace operations
func (workspace *Workspace) Stop() error {
	// Close the exit channel
	close(workspace.chanUpdateConfig)

	// Close all the modules
	for _, module := range workspace.modules {
		// TODO: Deal with the error in a better way
		_ = module.Close()
	}

	// Close the router
	workspace.router = nil

	// Close the config driver
	// TODO: Deal with the error in a better way
	_ = workspace.configDriver.Close()
	return nil
}
