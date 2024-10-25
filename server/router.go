package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YourTechBud/inferix/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (s *Server) createModuleRouter(module string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First acquire a read lock
		s.moduleLock.RLock()
		defer s.moduleLock.RUnlock()

		// Get the module
		router, ok := s.modules[module]
		if !ok {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Module not found", "invalid_module"))
			return
		}

		// Let the module handle the request
		router.ServeHTTP(w, r)
	})
}

func (s *Server) getGlobalConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the config
		config, err := s.configDriver.ReadAll(r.Context())
		if err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error reading configuration", "config_error"))
			return
		}

		// Write the config
		utils.WriteJSON(w, config)
	}
}

func (s *Server) createConfigRoutes(module string, resources []utils.ResourceConfiguration) http.Handler {
	router := chi.NewRouter()

	for _, resource := range resources {
		router.Post(fmt.Sprintf("/%s", resource.Path), s.configSetHandler(module, resource))
		router.Get(fmt.Sprintf("/%s", resource.Path), s.configGetHandler(module, resource))
		router.Delete(fmt.Sprintf("/%s/{id}", resource.Path), s.configDeleteHandler(module, resource))
	}

	return router
}

func (s *Server) configSetHandler(module string, cfg utils.ResourceConfiguration) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		// Read the body
		resource := cfg.New()
		if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Invalid request body", "invalid_request"))
			return
		}

		// Validate the resource
		if err := validate.Struct(resource); err != nil {
			errs := err.(validator.ValidationErrors)
			for _, e := range errs {
				log.Default().Println("field:", e.Field(), "error:", e.Tag())
			}
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Invalid request body", "invalid_request"))
			return
		}

		// Write the resource to the config store
		switch cfg.Type {
		case utils.ConfigResourceType_Object:
			if err := s.configDriver.SetInObject(r.Context(), module, cfg.Path, resource.GetID(), resource); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error writing configuration: %s", err), "config_error"))
				return
			}
		case utils.ConfigResourceType_Array:
			if err := s.configDriver.SetInArray(r.Context(), module, cfg.Path, resource.GetID(), resource); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error writing configuration: %s", err), "config_error"))
				return
			}
		}

		// Update the modules
		// TODO: Add support to revert back if the module could not be loaded for whatever reason
		s.updateConfig()

		// Write the response
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) configGetHandler(module string, cfg utils.ResourceConfiguration) http.HandlerFunc {
	type response struct {
		Resources json.RawMessage `json:"resources"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path
		path := cfg.Path

		// Get the resources
		resources, err := s.configDriver.Get(r.Context(), module, path)
		if err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error reading configuration", "config_error"))
			return
		}

		// Default to an empty object or array
		if resources == nil {
			switch cfg.Type {
			case utils.ConfigResourceType_Object:
				resources = json.RawMessage("{}")
			case utils.ConfigResourceType_Array:
				resources = json.RawMessage("[]")
			}
		}

		// Write the resources
		utils.WriteJSON(w, response{Resources: resources})
	}
}

func (s *Server) configDeleteHandler(module string, cfg utils.ResourceConfiguration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path and id
		path := cfg.Path

		// Delete the resource
		switch cfg.Type {
		case utils.ConfigResourceType_Object:
			if err := s.configDriver.DeleteFromObject(r.Context(), module, path, chi.URLParam(r, "id")); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error deleting configuration", "config_error"))
				return
			}
		case utils.ConfigResourceType_Array:
			if err := s.configDriver.DeleteFromArray(r.Context(), module, path, chi.URLParam(r, "id")); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error deleting configuration", "config_error"))
				return
			}
		}

		// Update the modules
		// TODO: Add support to revert back if the module could not be loaded for whatever reason
		s.updateConfig()

		// Write the response
		w.WriteHeader(http.StatusNoContent)
	}
}
