package utils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// ReadYAMLFile reads a file and unmarshals it as YAML
func ReadYAMLFile(path string, vPtr any) error {
	// Open a file descriptor
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read file contents
	yamlContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Convert YAML to JSON because we are fans of JSON
	jsonContent, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		return err
	}

	// Unmarshal the json content into the provided structure
	err = json.Unmarshal(jsonContent, vPtr)
	if err != nil {
		return err
	}

	return nil
}

// WriteYAMLFile writes a structure to a file as YAML
func WriteYAMLFile(path string, v any) error {
	// Convert the structure to JSON
	jsonContent, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Convert JSON to YAML
	yamlContent, err := yaml.JSONToYAML(jsonContent)
	if err != nil {
		return err
	}

	// Open a file descriptor
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the YAML content to the file
	_, err = file.Write(yamlContent)
	if err != nil {
		return err
	}

	return nil
}

// CreateDirIfNotExists creates a directory if it doesn't exist
func CreateDirIfNotExists(path string) error {
	// Extract the directory path
	dirPath := filepath.Dir(path)

	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// DeleteDirectory deletes a directory along with all its contents
func DeleteDirectory(path string) error {
	// Check if the directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// Remove the directory
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}

// CheckIfFileExists checks if a file exists
func CheckIfFileExists(path string) bool {
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
