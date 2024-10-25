package config

// Config is a struct to configure the LLM module
type Config struct {
	Models   []ModelConfig   `json:"models" yaml:"models"`
	Backends []BackendConfig `json:"backends" yaml:"backends"`
}
