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
			name: "Single Array Element",
			elements: sqliteConfigElements{
				{
					Module:    "module1",
					Path:      "path1",
					ElementID: "element1",
					Config:    `{"id": "element1", "key": "value"}`,
				},
			},
			expected: map[string]any{
				"module1": map[string]any{
					"path1": []any{
						json.RawMessage(`{"id": "element1", "key": "value"}`),
					},
				},
			},
		}, {
			name: "Multiple Elements",
			elements: sqliteConfigElements{
				{
					Module:    "module1",
					Path:      "path1",
					ElementID: "element1",
					Config:    `{"id": "element1", "key": "value1"}`,
				},
				{
					Module:    "module1",
					Path:      "path1",
					ElementID: "element2",
					Config:    `{"id": "element2", "key": "value2"}`,
				},
				{
					Module:    "module1",
					Path:      "path2",
					ElementID: "element1",
					Config:    `{"id": "element1", "key": "value3"}`,
				},
				{
					Module:    "module1",
					Path:      "path2",
					ElementID: "element2",
					Config:    `{"id": "element2", "key": "value4"}`,
				},
				{
					Module:    "module2",
					Path:      "path1",
					ElementID: "element1",
					Config:    `{"id": "element1", "key": "value5"}`,
				},
			},
			expected: map[string]any{
				"module1": map[string]any{
					"path1": []any{
						json.RawMessage(`{"id": "element1", "key": "value1"}`),
						json.RawMessage(`{"id": "element2", "key": "value2"}`),
					},
					"path2": []any{
						json.RawMessage(`{"id": "element1", "key": "value3"}`),
						json.RawMessage(`{"id": "element2", "key": "value4"}`),
					},
				},
				"module2": map[string]any{
					"path1": []any{
						json.RawMessage(`{"id": "element1", "key": "value5"}`),
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
