package backends

import "encoding/json"

// Config is a struct for backend configuration
type Config struct {
	Name        string          `json:"name"`
	BackendType string          `json:"type"`
	Config      json.RawMessage `json:"config"`
}
