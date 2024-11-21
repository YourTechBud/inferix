package server

import (
	"context"

	"github.com/YourTechBud/inferix/utils"
	"github.com/valyala/fastjson"
)

func (workspace *Workspace) configureModules(ctx context.Context) error {
	// Get the config
	cfg, err := workspace.configDriver.ReadAll(ctx)
	if err != nil {
		return err
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(cfg)
	if err != nil {
		return err
	}

	// Create all the modules
	modules := make([]utils.Module, len(modulesList))
	for i, moduleInfo := range modulesList {
		name := moduleInfo.Name
		workspaceCtx := utils.NewWorkspaceContext(workspace.tenant, workspace.workspace, name, workspace.configDriver)
		moduleConfig := v.Get(name).MarshalTo(nil)
		module, err := moduleInfo.New(workspaceCtx, moduleConfig)
		if err != nil {
			return err
		}

		modules[i] = module
	}

	// Don't forget to lock the workspace
	workspace.lock.Lock()
	defer workspace.lock.Unlock()

	// Close the old modules
	for _, module := range workspace.modules {
		// TODO: Deal with the error in a better way
		_ = module.Close()
	}

	// Set the new modules
	workspace.modules = modules

	// Start the router
	workspace.intializeRouter()

	return nil
}

func (workspace *Workspace) updateConfig() {
	workspace.chanUpdateConfig <- struct{}{}
}

func (workspace *Workspace) configUpdater() {
	for range workspace.chanUpdateConfig {
		// TODO: Add support to only configure the module which have been affected
		if err := workspace.configureModules(context.TODO()); err != nil {
			// TODO: Deal with the error in a better way
			panic(err)
		}
	}
}
