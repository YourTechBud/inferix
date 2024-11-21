package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TODO: This file shouldn't be in the utils package.

type (
	// WorkspaceContext is passed as context to every module during initialization
	WorkspaceContext struct {
		Storage *WorkspaceStorage
	}

	// WorkspaceStorage is a struct for modules to interact with the database
	WorkspaceStorage struct {
		module    string
		tenant    string
		workspace string

		configDriver ConfigDriver
	}

	// RequestContext is a struct for the context passed to modules during the lifecycle hooks
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
		New              func(ctx *WorkspaceContext, config json.RawMessage) (Module, error)
		GetResourcesInfo func() []ResourceInfo
	}

	// ResourceInfo is a struct for a configuration resource
	ResourceInfo struct {
		Path            string
		ProtectedFields []string
		New             func() Resource
	}

	// Lifecycle Hooks for resources

	ResourceProvisioner interface {
		Provision(ctx *RequestContext) (returningValue, metadata any, err error)
	}
	ResourceUpdater interface {
		Update(ctx *RequestContext, oldValue, oldMetadata any) (returningValue, metadata any, err error)
	}
	ResourceValidator interface {
		Validate() error
	}

	// The resource itself

	ResourceObject struct {
		Config    json.RawMessage `json:"config"`
		Metadata  json.RawMessage `json:"metadata"`
		CreatedAt int64           `json:"created_at"`
		UpdatedAt int64           `json:"updated_at"`
	}

	Resource interface {
		GetID() string
	}

	// RequestContextKeyType is the type for the server context key
	RequestContextKeyType string

	// ConfigDriverType is the type of configuration driver to use
	ConfigDriverType string

	// ConfigDriver is an interface for managing configuration
	ConfigDriver interface {
		SetResource(ctx context.Context, module, path, id string, element, metadata any) error
		SetResourceMetadata(ctx context.Context, module, path, id string, metadata any) error
		ReadAll(ctx context.Context) (json.RawMessage, error)
		GetAllResources(ctx context.Context, module, path string) ([]*ResourceObject, error)
		GetResource(ctx context.Context, module, path, id string) (*ResourceObject, error)
		DeleteResource(ctx context.Context, module, path, id string) error
		CheckIfResourceExists(ctx context.Context, module, path, id string) (bool, error)
		Close() error
	}
)

const (
	RequestContextKey RequestContextKeyType = "request_context"

	ConfigDriverType_File   ConfigDriverType = "file"
	ConfigDriverType_LibSQL ConfigDriverType = "libsql"
)

func NewWorkspaceContext(tenant, workspace, module string, configdriver ConfigDriver) *WorkspaceContext {
	return &WorkspaceContext{
		Storage: &WorkspaceStorage{
			tenant:       tenant,
			workspace:    workspace,
			module:       module,
			configDriver: configdriver,
		},
	}
}

func (s *WorkspaceStorage) UpdateConfigMetadata(ctx context.Context, path, elementID string, data any) error {
	return s.configDriver.SetResourceMetadata(ctx, s.module, path, elementID, data)
}

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

func GetRequestContext(r *http.Request) *RequestContext {
	return r.Context().Value(RequestContextKey).(*RequestContext)
}
