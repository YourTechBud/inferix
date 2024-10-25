package server

import (
	"context"
	"fmt"

	"github.com/YourTechBud/inferix/modules/llm"
	"github.com/valyala/fastjson"
)

func initialiseConfigDriver(opts Options) (ConfigDriver, error) {
	switch opts.ConfigDriver {
	case ConfigDriverType_File:
		return NewFileConfigDriver(opts)
	case ConfigDriverType_LibSQL:
		return NewLibSQLConfigDriver(opts)
	default:
		return nil, fmt.Errorf("unknown storage driver: %s", opts.ConfigDriver)
	}
}

func (s *Server) configureModules(ctx context.Context) error {
	// Get the config
	cfg, err := s.configDriver.ReadAll(ctx)
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
	s.moduleLock.Lock()
	s.modules["llm"] = llm.Routes()
	s.moduleLock.Unlock()

	return nil
}

func (s *Server) updateConfig() {
	s.chanUpdateConfig <- struct{}{}
}

func (s *Server) configUpdater() {
	for range s.chanUpdateConfig {
		// TODO: Add support to only configure the module which have been affected
		if err := s.configureModules(context.TODO()); err != nil {
			// TODO: Deal with the error in a better way
			panic(err)
		}
	}
}
