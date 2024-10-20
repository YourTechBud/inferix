package models

import "github.com/YourTechBud/inferix/modules/llm/types"

// Config represents the configuration for a model
type Config struct {
	Name           string        `json:"name"`
	Aliases        []string      `json:"aliases,omitempty"`
	Driver         string        `json:"driver"`
	TargetName     string        `json:"target_name,omitempty"`
	DefaultOptions *ModelOptions `json:"default_options,omitempty"`
}

// GetName returns the model's name
func (m *Config) GetName() string {
	return m.Name
}

// GetTargetName returns the target name if present, otherwise the model's name
func (m *Config) GetTargetName() string {
	if m.TargetName != "" {
		return m.TargetName
	}
	return m.Name
}

// ModelOptions represents options for configuring a model
type ModelOptions struct {
	TopP        *float64 `json:"top_p,omitempty"`
	TopK        *int32   `json:"top_k,omitempty"`
	NumCtx      *int32   `json:"num_ctx,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
}

// DefaultModelOptions returns default options for the model
func DefaultModelOptions() *ModelOptions {
	temperature := float64(0.2)
	return &ModelOptions{
		Temperature: &temperature,
	}
}

// MergeOptions merges the given options with the default options. It modifies the given options directly.
func (m *Config) MergeOptions(opts *types.InferenceOptions) {
	if opts.NumCtx == nil {
		opts.NumCtx = m.DefaultOptions.NumCtx
	}
	if opts.TopP == nil {
		opts.TopP = m.DefaultOptions.TopP
	}
	if opts.TopK == nil {
		opts.TopK = m.DefaultOptions.TopK
	}
	if opts.Temperature == nil {
		opts.Temperature = m.DefaultOptions.Temperature
	}
}
