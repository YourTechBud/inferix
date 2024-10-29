package server

import (
	"context"
	"encoding/json"
)

type (
	// The main configuration struct for the server
	Options struct {
		ConfigDriver      ConfigDriverType
		ConfigPath        string
		DefaultConfigPath string
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
	}
)

const (
	ConfigDriverType_File   ConfigDriverType = "file"
	ConfigDriverType_LibSQL ConfigDriverType = "libsql"
)
