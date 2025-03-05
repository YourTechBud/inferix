package config

import "github.com/YourTechBud/inferix/modules/llm/types"

// ModelConfig represents the configuration for a model
type ModelConfig struct {
	ID             string        `json:"id" validate:"required"`
	Aliases        []string      `json:"aliases,omitempty"`
	Backend        string        `json:"backend" validate:"required"`
	Target         string        `json:"target,omitempty"`
	DefaultOptions *ModelOptions `json:"default_options,omitempty"`
}

// GetName returns the model's name
func (m *ModelConfig) GetID() string {
	return m.ID
}

// GetTarget returns the target name if present, otherwise the model's name
func (m *ModelConfig) GetTarget() string {
	if m.Target != "" {
		return m.Target
	}
	return m.ID
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
	return &ModelOptions{}
}

// MergeOptions merges the given options with the default options. It modifies the given options directly.
func (m *ModelConfig) MergeOptions(opts *types.InferenceOptions) {
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
