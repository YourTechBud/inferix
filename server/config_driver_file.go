package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YourTechBud/inferix/utils"
)

// FileConfigDriver manages the configuration stored in a file
type FileConfigDriver struct {
	config   map[string]any
	filePath string
}

// NewFileConfigDriver creates a new FileConfigDriver
func NewFileConfigDriver(opts Options) (*FileConfigDriver, error) {
	// Throw an error if the default config path is set
	if opts.DefaultConfigPath != "" {
		return nil, fmt.Errorf("default config path is not supported for file config driver")
	}

	// Read the yaml file from the path provided
	var cfg map[string]any
	if err := utils.ReadYAMLFile(opts.ConfigPath, &cfg); err != nil {
		return nil, err
	}

	// Return the driver
	return &FileConfigDriver{
		config:   cfg,
		filePath: opts.ConfigPath,
	}, nil
}

func (f *FileConfigDriver) ReadAll(_ context.Context) (json.RawMessage, error) {
	data, _ := json.Marshal(f.config)
	return data, nil
}

func (f *FileConfigDriver) Get(_ context.Context, module, path string) (json.RawMessage, error) {
	// Get the value from the map
	value := utils.GetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), nil)

	if data, ok := value.(json.RawMessage); ok {
		return data, nil
	}

	return json.Marshal(value)
}

func (f *FileConfigDriver) SetInArray(_ context.Context, module, path, id string, element interface{}) error {
	// Convert the element to a map
	jsonData, _ := json.Marshal(element)
	var elementMap map[string]any
	_ = json.Unmarshal(jsonData, &elementMap)

	// Get the value from the map
	value := utils.GetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), []any{})
	if value == nil {
		return fmt.Errorf("path %s/%s not found", module, path)
	}

	// Convert the value to a slice
	slice, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("path %s/%s is not an array", module, path)
	}

	// Find the element in the slice
	found := false
	for i, v := range slice {
		if v.(map[string]any)["id"] == id {
			slice[i] = elementMap
			found = true
			break
		}
	}

	// If the element was not found, append it
	if !found {
		slice = append(slice, elementMap)
	}

	// Set the value in the map
	utils.SetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), slice)

	return utils.WriteYAMLFile(f.filePath, f.config)
}

func (f *FileConfigDriver) SetInObject(_ context.Context, module, path, id string, element interface{}) error {
	// Convert the element to a map
	jsonData, _ := json.Marshal(element)
	var elementMap map[string]any
	_ = json.Unmarshal(jsonData, &elementMap)

	// Get the value from the map
	value := utils.GetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), map[string]any{})
	if value == nil {
		return fmt.Errorf("path %s/%s not found", module, path)
	}

	// Convert the value to a map
	obj, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("path %s/%s is not an object", module, path)
	}

	// Set the value in the map
	obj[id] = elementMap

	// Set the value in the map
	utils.SetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), obj)

	return utils.WriteYAMLFile(f.filePath, f.config)
}

func (f *FileConfigDriver) DeleteFromArray(_ context.Context, module, path, id string) error {
	// Get the value from the map
	value := utils.GetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), nil)
	if value == nil {
		return fmt.Errorf("path %s/%s not found", module, path)
	}

	// Convert the value to a slice
	slice, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("path %s/%s is not an array", module, path)
	}

	// Find the element in the slice
	for i, v := range slice {
		if v.(map[string]any)["id"] == id {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	// Set the value in the map
	utils.SetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), slice)

	return utils.WriteYAMLFile(f.filePath, f.config)
}

func (f *FileConfigDriver) DeleteFromObject(_ context.Context, module, path, id string) error {
	// Get the value from the map
	value := utils.GetValueAtPath(f.config, fmt.Sprintf("%s/%s", module, path), nil)
	if value == nil {
		return fmt.Errorf("path %s/%s not found", module, path)
	}

	// Convert the value to a map
	obj, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("path %s/%s is not an object", module, path)
	}

	// Delete the value from the map
	delete(obj, id)

	return utils.WriteYAMLFile(f.filePath, f.config)
}

var _ ConfigDriver = (*FileConfigDriver)(nil)
