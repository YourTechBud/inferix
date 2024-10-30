package server

import (
	"context"

	"github.com/YourTechBud/inferix/modules/llm"
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
	llmConfig := v.Get("llm").MarshalTo(nil)
	llm, err := llm.New(llmConfig)
	if err != nil {
		return err
	}

	// Add the module
	workspace.lock.Lock()
	workspace.modules["llm"] = llm
	workspace.lock.Unlock()

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
