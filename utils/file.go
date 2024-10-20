package utils

import (
	"io"
	"os"

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
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Unmarshal the YAML content into the provided structure
	err = yaml.Unmarshal(content, vPtr)
	if err != nil {
		return err
	}

	return nil
}
