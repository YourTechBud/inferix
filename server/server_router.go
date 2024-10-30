package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YourTechBud/inferix/utils"
	"github.com/go-chi/chi/v5"
)

func (s *Server) workspaceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the `X-Inferix-Tenant` header
		tenant := r.Header.Get("X-Inferix-Tenant")
		if tenant == "" {
			tenant = "default"
		}

		// Read the `X-Inferix-Workspace` header
		workspace := r.Header.Get("X-Inferix-Workspace")
		if workspace == "" {
			workspace = "default"
		}

		serverContext := ServerContext{
			tenant:    tenant,
			workspace: workspace,
		}

		r = r.WithContext(context.WithValue(r.Context(), ServerContextKey, serverContext))

		next.ServeHTTP(w, r)
	})
}

func (s *Server) router() http.Handler {
	router := chi.NewRouter()

	// Add the workspace middleware
	router.Use(s.workspaceMiddleware)

	// Setup the workspace management routes
	router.Post("/inferix/v1/workspace", s.handleCreateWorkspace())
	router.Delete("/inferix/v1/workspace/{workspace}", s.handleDeleteWorkspace())

	// Setup the workspace routes
	router.Mount("/inferix/v1", s.handleWorkspaceRoutes())

	return router
}

func (s *Server) handleCreateWorkspace() http.HandlerFunc {
	type createWorkspaceRequest struct {
		Workspace string `json:"workspace"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the server context
		serverContext := r.Context().Value(ServerContextKey).(ServerContext)

		// Parse the request
		var req createWorkspaceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Invalid request", "invalid_request"))
			return
		}

		// Create the workspace
		if err := s.NewWorkspace(r.Context(), serverContext.tenant, req.Workspace); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error creating workspace - %s", err), "workspace_error"))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) handleDeleteWorkspace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the server context
		serverContext := r.Context().Value(ServerContextKey).(ServerContext)

		// Get the workspace
		workspace := chi.URLParam(r, "workspace")

		// Remove the workspace
		if err := s.RemoveWorkspace(serverContext.tenant, workspace); err != nil {
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
		serverContext := r.Context().Value(ServerContextKey).(ServerContext)

		// Get the workspace
		workspaceKey := getWorkspaceKey(serverContext.tenant, serverContext.workspace)

		// Get the workspace
		workspace, ok := s.workspaces[workspaceKey]
		if !ok {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Workspace not found", "invalid_workspace"))
			return
		}

		// Serve the request
		workspace.router.ServeHTTP(w, r)
	})
}
