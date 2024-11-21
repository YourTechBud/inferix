package security

import (
	"net/http"
	"strings"

	"github.com/YourTechBud/inferix/utils"
	"github.com/go-chi/chi/v5"
)

func (module *Module) Middlewares() []utils.HTTPMiddleware {
	return []utils.HTTPMiddleware{module.initializeMiddleware()}
}

func (module *Module) Routes() chi.Router {
	// We got no routes
	return nil
}

func (module *Module) initializeMiddleware() utils.HTTPMiddleware {
	return utils.HTTPMiddleware{
		RouteTypes: []utils.HTTPRouteType{utils.HTTPRouteType_API},
		Handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// First check if the request has already been authenticated
				if utils.GetRequestContext(r).Authenticated() {
					next.ServeHTTP(w, r)
					return
				}

				// Read the apiKey header
				apiKey := r.Header.Get("Authorization")

				// Strip the "Bearer " prefix
				apiKey = strings.TrimPrefix(apiKey, "Bearer ")

				// Validate the API key
				keyID, err := module.validate(apiKey)
				if err != nil {
					utils.WriteJSONError(w, err)
					return
				}

				// Update the last used time
				go func(id string) {
					module.workspaceStorage.UpdateConfigMetadata(r.Context(), "keys", id, Metadata{utils.CurrentTime()})
				}(keyID)

				next.ServeHTTP(w, r)
			})
		},
	}
}
