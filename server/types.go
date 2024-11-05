package server

import (
	"context"
	"encoding/json"

	"github.com/go-chi/chi/v5"
)

type (
	// The main configuration struct for the server
	Options struct {
		ConfigDriver           ConfigDriverType `mapstructure:"config-driver"`
		ConfigPath             string           `mapstructure:"config-path"`
		DefaultConfigPath      string           `mapstructure:"default-config-path"`
		CreateDefaultWorkspace bool             `mapstructure:"create-default-workspace"`
	}

	// ConfigDriverType is the type of configuration driver to use
	ConfigDriverType string

	// ConfigDriver is an interface for managing configuration
	ConfigDriver interface {
		SetInArray(ctx context.Context, module, path, id string, element interface{}) error
		SetInObject(ctx context.Context, module, path, id string, element interface{}) error
		ReadAll(ctx context.Context) (json.RawMessage, error)
		Get(ctx context.Context, module, path string) (json.RawMessage, error)
		DeleteFromArray(ctx context.Context, module, path, id string) error
		DeleteFromObject(ctx context.Context, module, path, id string) error
		CheckIfResourceExists(ctx context.Context, module, path, id string) (bool, error)
		Close() error
	}

	// Module describes the methods that a module must implement
	Module interface {
		Routes() chi.Router
		Close() error
	}

	// ServerContext holds the context for the server
	ServerContext struct {
		tenant, workspace string
	}

	// ServerContextKeyType is the type for the server context key
	ServerContextKeyType string
)

const (
	ServerContextKey ServerContextKeyType = "server_context"

	ConfigDriverType_File   ConfigDriverType = "file"
	ConfigDriverType_LibSQL ConfigDriverType = "libsql"
)
