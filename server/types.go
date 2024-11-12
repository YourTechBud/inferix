package server

import (
	"context"
	"encoding/json"
)

type (
	// The main configuration struct for the server
	Options struct {
		ConfigDriver           ConfigDriverType `mapstructure:"config-driver"`
		ConfigPath             string           `mapstructure:"config-path"`
		DefaultConfigPath      string           `mapstructure:"default-config-path"`
		CreateDefaultWorkspace bool             `mapstructure:"create-default-workspace"`
		AuthOptions            AuthOptions      `mapstructure:"auth"`
	}

	// AuthOptions describes the configuration for authentication
	AuthOptions struct {
		Enabled bool
		User    string
		Pass    string
	}

	// ConfigDriverType is the type of configuration driver to use
	ConfigDriverType string

	// ConfigDriver is an interface for managing configuration
	ConfigDriver interface {
		SetInArray(ctx context.Context, module, path, id string, element interface{}) error
		SetInObject(ctx context.Context, module, path, id string, element interface{}) error
		ReadAll(ctx context.Context) (json.RawMessage, error)
		GetAllResources(ctx context.Context, module, path string) (json.RawMessage, error)
		GetResource(ctx context.Context, module, path, id string) (json.RawMessage, error)
		DeleteFromArray(ctx context.Context, module, path, id string) error
		DeleteFromObject(ctx context.Context, module, path, id string) error
		CheckIfResourceExists(ctx context.Context, module, path, id string) (bool, error)
		Close() error
	}
)

const (
	ConfigDriverType_File   ConfigDriverType = "file"
	ConfigDriverType_LibSQL ConfigDriverType = "libsql"
)
