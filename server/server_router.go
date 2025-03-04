package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/YourTechBud/inferix/modules/security"
	"github.com/YourTechBud/inferix/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) middlewareAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First get the request context
		requestContext := utils.GetRequestContext(r)

		// Check if authentication is enabled or not
		if !s.options.AuthOptions.Enabled {
			// Simply mark the request as authenticated and move on.
			requestContext.SetAuthenticated(true)
			next.ServeHTTP(w, r)
			return
		}

		// Check if the Authorization header contains basic auth
		username, password, ok := r.BasicAuth()
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// Check if the username and password are correct
		if username == s.options.AuthOptions.User && password == s.options.AuthOptions.Pass {
			requestContext.SetAuthenticated(true)
			next.ServeHTTP(w, r)
			return
		}

		// Return a 401
		utils.WriteJSONError(w, utils.NewStandardError(http.StatusUnauthorized, "Invalid credentials", "invalid_credentials"))
	})
}

func (s *Server) middlewareServerContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create the tenant and workspace variables
		var tenant, workspace string

		// First check if we have received an apikey. API keys take top priority
		apiKey := r.Header.Get("Authorization")
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")
		if apiKey != "" && strings.HasPrefix(apiKey, "ik:") {
			var err error
			tenant, workspace, _, _, err = security.GetKeySegments(apiKey)
			if err != nil {
				utils.WriteJSONError(w, err)
				return
			}
		}

		// Check the headers only if we didn't get the tenant and workspace from the apikey
		if tenant == "" || workspace == "" {
			// Read the `X-Inferix-Tenant` header
			if t := r.Header.Get("X-Inferix-Tenant"); t != "" {
				tenant = t
			}

			// Read the `X-Inferix-Workspace` header
			if w := r.Header.Get("X-Inferix-Workspace"); w != "" {
				workspace = w
			}
		}

		// If the tenant or workspace is empty, set them to default
		if tenant == "" {
			tenant = "default"
		}
		if workspace == "" {
			workspace = "default"
		}

		serverContext := utils.NewRequestContext(tenant, workspace)
		r = r.WithContext(context.WithValue(r.Context(), utils.RequestContextKey, serverContext))

		next.ServeHTTP(w, r)
	})
}

func (s *Server) middlewareLoadWorkspace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the server context
		serverContext := utils.GetRequestContext(r)

		// Try loading the workspace into memory
		if err := s.LoadWorkspace(r.Context(), serverContext.Tenant(), serverContext.Workspace()); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, fmt.Sprintf("Error loading workspace - %s", err), "workspace_error"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) router() http.Handler {
	router := chi.NewRouter()

	// Setup CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	// Setup the server context middleware
	router.Use(s.middlewareServerContext, s.middlewareAuthentication)

	// Setup the workspace management routes
	router.Post("/inferix/v1/workspace", s.handleCreateWorkspace())
	router.Delete("/inferix/v1/workspace/{workspace}", s.handleDeleteWorkspace())

	// Setup the workspace routes
	router.Mount("/inferix/v1", s.middlewareLoadWorkspace(s.handleWorkspaceRoutes()))

	return router
}

func (s *Server) handleCreateWorkspace() http.HandlerFunc {
	type createWorkspaceRequest struct {
		Workspace string `json:"workspace"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the server context
		serverContext := utils.GetRequestContext(r)

		// Parse the request
		var req createWorkspaceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Invalid request", "invalid_request"))
			return
		}

		// Create the workspace
		if err := s.NewWorkspace(r.Context(), serverContext.Tenant(), req.Workspace); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error creating workspace - %s", err), "workspace_error"))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) handleDeleteWorkspace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the server context
		serverContext := utils.GetRequestContext(r)

		// Get the workspace
		workspace := chi.URLParam(r, "workspace")

		// Remove the workspace
		if err := s.RemoveWorkspace(r.Context(), serverContext.Tenant(), workspace); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error removing workspace", "workspace_error"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleWorkspaceRoutes() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.lock.RLock()
		defer s.lock.RUnlock()

		// First get the server context
		serverContext := utils.GetRequestContext(r)

		// Get the workspace
		workspaceKey := getWorkspaceKey(serverContext.Tenant(), serverContext.Workspace())

		// Get the workspace
		workspace := s.workspaces[workspaceKey]

		// Serve the request
		workspace.router.ServeHTTP(w, r)
	})
}
