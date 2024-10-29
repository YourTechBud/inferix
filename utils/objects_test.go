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
			if !assert.Equal(t, tt.expected, tt.obj) {
				t.Errorf("expected %v, got %v", tt.expected, tt.obj)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	tests := []struct {
		name     string
		dst      map[string]any
		src      map[string]any
		expected map[string]any
	}{
		{
			name: "src nil",
			dst: map[string]any{
				"key1": "value1",
			},
			src: nil,
			expected: map[string]any{
				"key1": "value1",
			},
		},
		{
			name: "simple merge",
			dst: map[string]any{
				"key1": "value1",
			},
			src: map[string]any{
				"key2": "value2",
			},
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "nested merge",
			dst: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
				},
			},
			src: map[string]any{
				"key1": map[string]any{
					"key3": "value3",
				},
			},
			expected: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
					"key3": "value3",
				},
			},
		},
		{
			name: "overwrite value",
			dst: map[string]any{
				"key1": "value1",
			},
			src: map[string]any{
				"key1": "newValue1",
			},
			expected: map[string]any{
				"key1": "newValue1",
			},
		},
		{
			name: "merge with empty dst",
			dst:  map[string]any{},
			src: map[string]any{
				"key1": "value1",
			},
			expected: map[string]any{
				"key1": "value1",
			},
		},
		{
			name: "merge with empty src",
			dst: map[string]any{
				"key1": "value1",
			},
			src: map[string]any{},
			expected: map[string]any{
				"key1": "value1",
			},
		},
		{
			name: "complex nested merge",
			dst: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
					"key3": map[string]any{
						"key4": "value4",
					},
				},
			},
			src: map[string]any{
				"key1": map[string]any{
					"key3": map[string]any{
						"key5": "value5",
					},
					"key6": "value6",
				},
			},
			expected: map[string]any{
				"key1": map[string]any{
					"key2": "value2",
					"key3": map[string]any{
						"key4": "value4",
						"key5": "value5",
					},
					"key6": "value6",
				},
			},
		},
		{
			name: "array appending",
			dst: map[string]any{
				"key1": []any{"value1"},
			},
			src: map[string]any{
				"key1": []any{"value2"},
			},
			expected: map[string]any{
				"key1": []any{"value1", "value2"},
			},
		},
		{
			name: "array appending with nested map",
			dst: map[string]any{
				"key1": map[string]any{
					"key2": []any{"value1"},
				},
			},
			src: map[string]any{
				"key1": map[string]any{
					"key2": []any{"value2"},
				},
			},
			expected: map[string]any{
				"key1": map[string]any{
					"key2": []any{"value1", "value2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeMaps(tt.dst, tt.src)
			if !assert.Equal(t, tt.expected, tt.dst) {
				t.Errorf("expected %v, got %v", tt.expected, tt.dst)
			}
		})
	}
}
