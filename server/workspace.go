package server

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/go-libsql"

	"github.com/YourTechBud/inferix/utils"
)

// Workspace is a map of modules
type Workspace struct {
	lock    sync.RWMutex
	modules []utils.Module

	// Metadata
	tenant    string
	workspace string

	// Config driver for the workspace
	db           *sqlx.DB
	configDriver utils.ConfigDriver

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

	// Check if the workspace already exists
	// TODO: How do we handle this for configdrivers which are not file based?
	configPath := getWorkspaceDBPath(s.options.StorageDirectory, tenant, workspace)
	if utils.CheckIfFileExists(configPath) {
		return fmt.Errorf("workspace already exists")
	}

	// Load the workspace into memory
	if err := s.loadWorkspace(ctx, tenant, workspace); err != nil {
		return err
	}

	return nil
}

// RemoveWorkspace removes a workspace
func (s *Server) RemoveWorkspace(ctx context.Context, tenant, workspace string) error {
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

	// Delete all the workspace files
	// TODO: How do we handle this for configdrivers which are not file based?
	configPath := getWorkspaceDir(s.options.StorageDirectory, tenant, workspace)
	_ = utils.DeleteDirectory(configPath)

	return nil
}

func (s *Server) LoadWorkspace(ctx context.Context, tenant, workspace string) error {
	// Get the workspace key
	workspaceKey := getWorkspaceKey(tenant, workspace)

	// Check if the workspace exists in memory
	s.lock.RLock()
	if _, p := s.workspaces[workspaceKey]; p {
		s.lock.RUnlock()
		return nil
	}
	s.lock.RUnlock()

	// Check if the workspace exists in the config driver
	// TODO: How do we handle this for configdrivers which are not file based?
	configPath := getWorkspaceDBPath(s.options.StorageDirectory, tenant, workspace)
	if !utils.CheckIfFileExists(configPath) {
		return fmt.Errorf("workspace does not exist")
	}

	// Acquire the server lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// Load the workspace into memory
	return s.loadWorkspace(ctx, tenant, workspace)
}

// loadWorkspace loads a workspace into memory. Always make sure that this function is called
// after acquiring the server lock and after making sure that the workspace does exist.
func (s *Server) loadWorkspace(ctx context.Context, tenant, workspace string) error {
	log.Default().Printf("Loading workspace %s/%s", tenant, workspace)

	// Get the workspace key
	workspaceKey := getWorkspaceKey(tenant, workspace)

	// Setup the path for the workspace
	configPath := getWorkspaceDBPath(s.options.StorageDirectory, tenant, workspace)

	// First create the containing directory if it doesn't exist
	if err := utils.CreateDirIfNotExists(configPath); err != nil {
		return err
	}

	// Open the database
	db, err := sqlx.Open(s.options.StorageDriver, fmt.Sprintf("file://%s", configPath))
	if err != nil {
		return err
	}

	// Create a config driver for the workspace
	configDriver, err := NewLibSQLConfigDriver(db, s.defaultConfig)
	if err != nil {
		return err
	}

	// Create and return the workspace
	w := &Workspace{
		// Create the modules map
		modules: make([]utils.Module, 0),

		// Metadata
		tenant:    tenant,
		workspace: workspace,

		// Set the config driver
		db:           db,
		configDriver: configDriver,

		// Create the channels
		chanUpdateConfig: make(chan struct{}, 5),
	}

	// Start the workspace
	if err := w.Start(ctx); err != nil {
		return err
	}

	// Add the workspace to the server
	s.workspaces[workspaceKey] = w

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
	_ = workspace.db.Close()
	return nil
}
