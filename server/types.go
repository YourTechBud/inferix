package server

import "github.com/YourTechBud/inferix/utils"

type (
	// The main configuration struct for the server
	Options struct {
		ConfigDriver           utils.ConfigDriverType `mapstructure:"config-driver"`
		ConfigPath             string                 `mapstructure:"config-path"`
		DefaultConfigPath      string                 `mapstructure:"default-config-path"`
		CreateDefaultWorkspace bool                   `mapstructure:"create-default-workspace"`
		AuthOptions            AuthOptions            `mapstructure:"auth"`
	}

	// AuthOptions describes the configuration for authentication
	AuthOptions struct {
		Enabled bool
		User    string
		Pass    string
	}
)
