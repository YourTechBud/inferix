package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fastjson"

	"github.com/YourTechBud/inferix/utils"
)

func (workspace *Workspace) intializeRouter() {
	// Prepare all the routers
	router := chi.NewRouter()

	configRouter := chi.NewRouter()
	apiRouter := chi.NewRouter()

	// Add the authentication middleware for the config routes
	configRouter.Use(utils.ValidateAuthencation)

	// First setup all module middlewares
	for _, module := range workspace.modules {
		for _, middleware := range module.Middlewares() {
			for _, routeType := range middleware.RouteTypes {
				fmt.Println("Adding middleware for route type", routeType)
				switch routeType {
				case utils.HTTPRouteType_Config:
					configRouter.Use(middleware.Handler)
				case utils.HTTPRouteType_API:
					apiRouter.Use(middleware.Handler)
				}
			}
		}
	}

	// Setup module specific routes
	for i, moduleInfo := range modulesList {
		name := moduleInfo.Name

		// Setup the api routes if any
		if routes := workspace.modules[i].Routes(); routes != nil {
			apiRouter.Mount(fmt.Sprintf("/%s", name), workspace.createModuleRouter(workspace.modules[i].Routes()))
		}

		// Setup the config routes
		configRouter.Mount(fmt.Sprintf("/%s", name), workspace.createConfigRoutes(name, moduleInfo.GetResourcesInfo()))
	}

	// Setup the global config route
	configRouter.Get("/", workspace.getGlobalConfigHandler())

	// Mount the routers
	router.Mount("/config", configRouter)
	router.Mount("/", apiRouter)

	workspace.router = router
}

func (workspace *Workspace) createModuleRouter(moduleRouter chi.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First acquire a read lock
		workspace.lock.RLock()
		defer workspace.lock.RUnlock()

		// Let the module handle the request
		moduleRouter.ServeHTTP(w, r)
	})
}

func (workspace *Workspace) getGlobalConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the config
		config, err := workspace.configDriver.ReadAll(r.Context())
		if err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error reading configuration", "config_error"))
			return
		}

		// Remove the protected fields for each module
		var p fastjson.Parser
		parsedConfig, _ := p.Parse(string(config))
		for _, moduleInfo := range modulesList {
			// Get the resources declared by the module
			resources := moduleInfo.GetResourcesInfo()

			// Get the module configuration
			moduleValue := parsedConfig.Get(moduleInfo.Name)
			for _, resource := range resources {
				// Get the resource specific configuration
				resourceValue := moduleValue.Get(strings.Split(resource.Path, "/")...)
				removeFields(resourceValue, resource.ProtectedFields)
			}
		}

		// Write the config
		utils.WriteJSON(w, json.RawMessage(parsedConfig.MarshalTo(nil)))
	}
}

func (workspace *Workspace) createConfigRoutes(module string, resources []utils.ResourceInfo) http.Handler {
	router := chi.NewRouter()

	for _, resourceInfo := range resources {
		router.Post(fmt.Sprintf("/%s", resourceInfo.Path), workspace.configSetHandler(module, resourceInfo))
		router.Get(fmt.Sprintf("/%s", resourceInfo.Path), workspace.configGetHandler(module, resourceInfo))
		router.Delete(fmt.Sprintf("/%s/{id}", resourceInfo.Path), workspace.configDeleteHandler(module, resourceInfo))
	}

	return router
}

func (workspace *Workspace) configSetHandler(module string, resourceInfo utils.ResourceInfo) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		// Read the body
		resource := resourceInfo.New()
		if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Invalid request body", "invalid_request"))
			return
		}

		// Validate the json schema of the resource
		if err := validate.Struct(resource); err != nil {
			errs := err.(validator.ValidationErrors)
			for _, e := range errs {
				log.Default().Println("field:", e.Field(), "error:", e.Tag())
			}
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Invalid request body", "invalid_request"))
			return
		}

		// The value to return to the user.
		var returningValue any = nil

		// Check if the resource already doesResourceExists
		doesResourceExists, err := workspace.configDriver.CheckIfResourceExists(r.Context(), module, resourceInfo.Path, resource.GetID())
		if err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error checking if resource exists: %s", err), "config_error"))
			return
		}
		if !doesResourceExists {
			// Initialize the resource if it supports it
			if provisioner, ok := resource.(utils.ResourceProvisioner); ok {
				val, err := provisioner.Provision(r.Context().Value(utils.RequestContextKey).(*utils.RequestContext))
				if err != nil {
					utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, fmt.Sprintf("Error initializing resource: %s", err), "invalid_request"))
					return
				}

				returningValue = val
			}
		} else {
			// Update the resource if it supports it
			if updater, ok := resource.(utils.ResourceUpdater); ok {

				// Load the old value
				oldValue, err := workspace.configDriver.GetResource(r.Context(), module, resourceInfo.Path, resource.GetID())
				if err != nil {
					utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error reading configuration: %s", err), "config_error"))
					return
				}

				// Invoke the lifecycle hook
				oldResource := resourceInfo.New()
				_ = json.Unmarshal(oldValue, oldResource)
				val, err := updater.Update(r.Context().Value(utils.RequestContextKey).(*utils.RequestContext), oldResource)
				if err != nil {
					utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, fmt.Sprintf("Error updating resource: %s", err), "invalid_request"))
					return
				}

				returningValue = val
			}
		}

		// Validate the resource if it supports it
		if validator, ok := resource.(utils.ResourceValidator); ok {
			if err := validator.Validate(); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, fmt.Sprintf("Error validating resource: %s", err), "invalid_request"))
				return
			}
		}

		// Write the resource to the config store
		switch resourceInfo.Type {
		case utils.ConfigResourceType_Object:
			if err := workspace.configDriver.SetInObject(r.Context(), module, resourceInfo.Path, resource.GetID(), resource); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error writing configuration: %s", err), "config_error"))
				return
			}
		case utils.ConfigResourceType_Array:
			if err := workspace.configDriver.SetInArray(r.Context(), module, resourceInfo.Path, resource.GetID(), resource); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, fmt.Sprintf("Error writing configuration: %s", err), "config_error"))
				return
			}
		}

		// Update the modules
		// TODO: Add support to revert back if the module could not be loaded for whatever reason
		workspace.updateConfig()

		// Write the response
		if returningValue != nil {
			utils.WriteJSON(w, returningValue)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (workspace *Workspace) configGetHandler(module string, resourceInfo utils.ResourceInfo) http.HandlerFunc {
	type response struct {
		Resources json.RawMessage `json:"resources"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path
		path := resourceInfo.Path

		// Get the resources
		resources, err := workspace.configDriver.GetAllResources(r.Context(), module, path)
		if err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error reading configuration", "config_error"))
			return
		}

		// Default to an empty object or array
		if resources == nil {
			switch resourceInfo.Type {
			case utils.ConfigResourceType_Object:
				resources = json.RawMessage("{}")
			case utils.ConfigResourceType_Array:
				resources = json.RawMessage("[]")
			}
		}

		// Don't forget to remove the protected fields
		var p fastjson.Parser
		parsedResponse, _ := p.Parse(string(resources))
		removeFields(parsedResponse, resourceInfo.ProtectedFields)

		// Write the resources
		utils.WriteJSON(w, response{Resources: parsedResponse.MarshalTo(nil)})
	}
}

func (workspace *Workspace) configDeleteHandler(module string, resourceInfo utils.ResourceInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path and id
		path := resourceInfo.Path

		// Delete the resource
		switch resourceInfo.Type {
		case utils.ConfigResourceType_Object:
			if err := workspace.configDriver.DeleteFromObject(r.Context(), module, path, chi.URLParam(r, "id")); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error deleting configuration", "config_error"))
				return
			}
		case utils.ConfigResourceType_Array:
			if err := workspace.configDriver.DeleteFromArray(r.Context(), module, path, chi.URLParam(r, "id")); err != nil {
				utils.WriteJSONError(w, utils.NewStandardError(http.StatusInternalServerError, "Error deleting configuration", "config_error"))
				return
			}
		}

		// Update the modules
		// TODO: Add support to revert back if the module could not be loaded for whatever reason
		workspace.updateConfig()

		// Write the response
		w.WriteHeader(http.StatusNoContent)
	}
}
