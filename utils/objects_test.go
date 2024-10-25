package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValueAtPath(t *testing.T) {
	tests := []struct {
		name         string
		obj          map[string]any
		path         string
		defaultValue any
		expected     any
	}{
		{
			name: "simple path",
			obj: map[string]any{
				"key1": "value1",
			},
			path:         "key1",
			defaultValue: "default",
			expected:     "value1",
		},
		{
			name: "nested path",
			obj: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
				},
			},
			path:         "key1/key2",
			defaultValue: "default",
			expected:     "value2",
		},
		{
			name: "non-existent path",
			obj: map[string]any{
				"key1": "value1",
			},
			path:         "key2",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name: "partially existent path",
			obj: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
				},
			},
			path:         "key1/key3",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name: "empty path",
			obj: map[string]any{
				"key1": "value1",
			},
			path:         "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetValueAtPath(tt.obj, tt.path, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
func TestSetValueAtPath(t *testing.T) {
	tests := []struct {
		name     string
		obj      map[string]any
		path     string
		value    any
		expected map[string]any
	}{
		{
			name: "simple path",
			obj: map[string]any{
				"key1": "value1",
			},
			path:  "key2",
			value: "value2",
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "nested path",
			obj: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
				},
			},
			path:  "key1/key3",
			value: "value3",
			expected: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
					"key3": "value3",
				},
			},
		},
		{
			name: "overwrite existing value",
			obj: map[string]any{
				"key1": "value1",
			},
			path:  "key1",
			value: "newValue1",
			expected: map[string]any{
				"key1": "newValue1",
			},
		},
		{
			name:  "create nested path",
			obj:   map[string]any{},
			path:  "key1/key2/key3",
			value: "value3",
			expected: map[string]any{
				"key1": map[string]any{
					"key2": map[string]any{
						"key3": "value3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetValueAtPath(tt.obj, tt.path, tt.value)
			if !assert.Equal(t, tt.obj, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, tt.obj)
			}
		})
	}
}
