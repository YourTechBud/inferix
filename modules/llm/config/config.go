package config

// Config is a struct to configure the LLM module
type Config struct {
	ModelConfigs   []ModelConfig   `json:"models" yaml:"models"`
	BackendConfigs []BackendConfig `json:"backends" yaml:"backends"`
	ModelSettings  []ModelSettings `json:"model-settings,omitempty" yaml:"model-settings,omitempty"`
}
