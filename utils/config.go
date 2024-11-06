package utils

import (
	"encoding/json"

	"github.com/go-chi/chi/v5"
)

type (
	// Module describes the methods that a module must implement
	Module interface {
		Routes() chi.Router
		Close() error
	}

	// ModuleInfo is a struct for module information
	ModuleInfo struct {
		New              func(config json.RawMessage) (Module, error)
		GetResourcesInfo func() []ResourceInfo
	}

	// ResourceInfo is a struct for a configuration resource
	ResourceInfo struct {
		Path            string
		Type            ConfigResourceType
		ProtectedFields []string
		New             func() Resource
	}

	// Lifecycle Hooks for resources

	ResourceProvisioner interface {
		Provision() error
	}
	ResourceValidator interface {
		Validate() error
	}

	// The resource itself

	Resource interface {
		GetID() string
	}

	// ConfigResourceType is the type of configuration resource
	ConfigResourceType string
)

const (
	ConfigResourceType_Object ConfigResourceType = "OBJECT"
	ConfigResourceType_Array  ConfigResourceType = "ARRAY"
)
