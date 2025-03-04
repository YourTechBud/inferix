package config

import (
	"errors"
	"regexp"

	"github.com/YourTechBud/inferix/utils"
)

// ModelSettings represents global settings that can be applied to models based on filters
type ModelSettings struct {
	// ID of the setting resource
	ID string `json:"id" yaml:"id" validate:"required"`

	// Options to be applied to matching models
	Options *ModelOptions `json:"options" yaml:"options"`

	// Filters specifies which models these options should apply to
	Filters []ModelFilter `json:"filters" yaml:"filters" validate:"required,min=1"`
}

// ModelFilter represents criteria for filtering models
type ModelFilter struct {
	// Regex pattern to match against model IDs
	Pattern string `json:"pattern" yaml:"pattern" validate:"required"`

	// Compiled regex pattern (not exported, used internally)
	compiledPattern *regexp.Regexp
}

// GetID returns the ID of the setting resource
func (s *ModelSettings) GetID() string {
	return s.ID
}

// Compile compiles the regex pattern in the filter
func (f *ModelFilter) Compile() error {
	pattern, err := regexp.Compile(f.Pattern)
	if err != nil {
		return err
	}
	f.compiledPattern = pattern
	return nil
}

// Matches checks if a model ID matches this filter
func (f *ModelFilter) Matches(modelID string) bool {
	if f.compiledPattern == nil {
		if err := f.Compile(); err != nil {
			panic(err) // This should never happen we validate the filters when they are added via the API
		}
	}
	return f.compiledPattern.MatchString(modelID)
}

// ApplySettings applies the model settings to a model config if it matches any of the filters
func (s *ModelSettings) ApplySettings(model *ModelConfig) {
	// Check if the model matches any of our filters
	matches := false
	for _, filter := range s.Filters {
		if filter.Matches(model.ID) {
			matches = true
			break
		}
	}

	if !matches {
		return
	}

	// If model doesn't have default options, initialize them
	if model.DefaultOptions == nil {
		model.DefaultOptions = DefaultModelOptions()
	}

	// Apply each option if it's set in the settings
	if s.Options.TopP != nil {
		model.DefaultOptions.TopP = s.Options.TopP
	}
	if s.Options.TopK != nil {
		model.DefaultOptions.TopK = s.Options.TopK
	}
	if s.Options.NumCtx != nil {
		model.DefaultOptions.NumCtx = s.Options.NumCtx
	}
	if s.Options.Temperature != nil {
		model.DefaultOptions.Temperature = s.Options.Temperature
	}
}

// Validate validates the model settings
func (s *ModelSettings) Validate() error {
	if len(s.Filters) == 0 {
		return errors.New("filters must contain at least one item")
	}
	for _, filter := range s.Filters {
		if err := filter.Compile(); err != nil {
			return err
		}
	}
	return nil
}

var _ utils.ResourceValidator = (*ModelSettings)(nil)
