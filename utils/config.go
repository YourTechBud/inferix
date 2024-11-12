package utils

import (
	"encoding/json"

	"github.com/go-chi/chi/v5"
)

// TODO: This file shouldn't be in the utils package.

type (
	// RequestContext is a struct for the context passed to modules during initialization
	RequestContext struct {
		tenant    string
		workspace string

		// Felds for authentication
		authenticated bool
	}

	// Module describes the methods that a module must implement
	Module interface {
		Routes() chi.Router
		Middlewares() []HTTPMiddleware
		Close() error
	}

	// ModuleInfo is a struct for module information
	ModuleInfo struct {
		Name             string
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
		Provision(ctx *RequestContext) (any, error)
	}
	ResourceUpdater interface {
		Update(ctx *RequestContext, oldValue any) (any, error)
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

	// RequestContextKeyType is the type for the server context key
	RequestContextKeyType string
)

const (
	RequestContextKey RequestContextKeyType = "request_context"

	ConfigResourceType_Object ConfigResourceType = "OBJECT"
	ConfigResourceType_Array  ConfigResourceType = "ARRAY"
)

// NewRequestContext creates a new context for a module
func NewRequestContext(tenant, workspace string) *RequestContext {
	return &RequestContext{
		tenant:    tenant,
		workspace: workspace,
	}
}

func (c *RequestContext) Tenant() string {
	return c.tenant
}

func (c *RequestContext) Workspace() string {
	return c.workspace
}

func (c *RequestContext) SetAuthenticated(authenticated bool) {
	c.authenticated = authenticated
}

func (c *RequestContext) Authenticated() bool {
	return c.authenticated
}
