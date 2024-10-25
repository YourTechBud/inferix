package server

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToMap(t *testing.T) {
	tests := []struct {
		name     string
		elements sqliteConfigElements
		expected map[string]any
	}{
		{
			name: "Single Object Element",
			elements: sqliteConfigElements{
				{
					Module:      "module1",
					Path:        "path1",
					ElementID:   "element1",
					ElementType: "OBJECT",
					Config:      `{"key": "value"}`,
				},
			},
			expected: map[string]any{
				"module1": map[string]any{
					"path1": map[string]any{
						"element1": json.RawMessage(`{"key": "value"}`),
					},
				},
			},
		},
		{
			name: "Single Array Element",
			elements: sqliteConfigElements{
				{
					Module:      "module1",
					Path:        "path1",
					ElementID:   "element1",
					ElementType: "ARRAY",
					Config:      `{"id": "element1", "key": "value"}`,
				},
			},
			expected: map[string]any{
				"module1": map[string]any{
					"path1": []any{
						json.RawMessage(`{"id": "element1", "key": "value"}`),
					},
				},
			},
		},
		{
			name: "Multiple Elements",
			elements: sqliteConfigElements{
				{
					Module:      "module1",
					Path:        "path1",
					ElementID:   "element1",
					ElementType: "OBJECT",
					Config:      `{"key": "value1"}`,
				},
				{
					Module:      "module1",
					Path:        "path1",
					ElementID:   "element2",
					ElementType: "OBJECT",
					Config:      `{"key": "value2"}`,
				},
				{
					Module:      "module1",
					Path:        "path2",
					ElementID:   "element1",
					ElementType: "ARRAY",
					Config:      `{"id": "element1", "key": "value3"}`,
				},
			},
			expected: map[string]any{
				"module1": map[string]any{
					"path1": map[string]any{
						"element1": json.RawMessage(`{"key": "value1"}`),
						"element2": json.RawMessage(`{"key": "value2"}`),
					},
					"path2": []any{
						json.RawMessage(`{"id": "element1", "key": "value3"}`),
					},
				},
			},
		},
		{
			name: "Multiple Object and Array Elements",
			elements: sqliteConfigElements{
				{
					Module:      "module1",
					Path:        "path1",
					ElementID:   "element1",
					ElementType: "OBJECT",
					Config:      `{"key": "value1"}`,
				},
				{
					Module:      "module1",
					Path:        "path1",
					ElementID:   "element2",
					ElementType: "OBJECT",
					Config:      `{"key": "value2"}`,
				},
				{
					Module:      "module1",
					Path:        "path2",
					ElementID:   "element1",
					ElementType: "ARRAY",
					Config:      `{"id": "element1", "key": "value3"}`,
				},
				{
					Module:      "module1",
					Path:        "path2",
					ElementID:   "element2",
					ElementType: "ARRAY",
					Config:      `{"id": "element2", "key": "value4"}`,
				},
				{
					Module:      "module2",
					Path:        "path1",
					ElementID:   "element1",
					ElementType: "OBJECT",
					Config:      `{"key": "value5"}`,
				},
			},
			expected: map[string]any{
				"module1": map[string]any{
					"path1": map[string]any{
						"element1": json.RawMessage(`{"key": "value1"}`),
						"element2": json.RawMessage(`{"key": "value2"}`),
					},
					"path2": []any{
						json.RawMessage(`{"id": "element1", "key": "value3"}`),
						json.RawMessage(`{"id": "element2", "key": "value4"}`),
					},
				},
				"module2": map[string]any{
					"path1": map[string]any{
						"element1": json.RawMessage(`{"key": "value5"}`),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.elements.convertToMap()
			assert.Equal(t, tt.expected, result)
		})
	}
}
