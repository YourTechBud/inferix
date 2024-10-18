package config

import "encoding/json"

// BackendConfig is a struct for backend configuration
type BackendConfig struct {
	Name        string          `json:"name"`
	BackendType string          `json:"type"`
	Config      json.RawMessage `json:"config"`
	Options     BackendOptions  `json:"options"`
}

// BackendOptions is a struct for backend options
type BackendOptions struct {
	InjectFnCallPrompt bool `json:"injectFnCallPrompt"`
}
