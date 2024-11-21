package server

type (
	// The main configuration struct for the server
	Options struct {
		StorageDriver          string      `mapstructure:"storage-driver"`
		StorageDirectory       string      `mapstructure:"storage-dir"`
		DefaultConfigPath      string      `mapstructure:"default-config-path"`
		CreateDefaultWorkspace bool        `mapstructure:"create-default-workspace"`
		AuthOptions            AuthOptions `mapstructure:"auth"`
	}

	// AuthOptions describes the configuration for authentication
	AuthOptions struct {
		Enabled bool
		User    string
		Pass    string
	}
)
