package apikeys

import (
	"net/http"
	"strings"

	"github.com/YourTechBud/inferix/utils"
	"github.com/go-chi/chi/v5"
)

func (module *Module) Middlewares() []utils.HTTPMiddleware {
	return module.middlewares
}

func (module *Module) Routes() chi.Router {
	// We got no routes
	return nil
}

func initializeMiddleware(config *Config) utils.HTTPMiddleware {
	return utils.HTTPMiddleware{
		RouteTypes: []utils.HTTPRouteType{utils.HTTPRouteType_API},
		Handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// First check if the request has already been authenticated
				if r.Context().Value(utils.RequestContextKey).(*utils.RequestContext).Authenticated() {
					next.ServeHTTP(w, r)
					return
				}

				// Read the apiKey header
				apiKey := r.Header.Get("Authorization")

				// Strip the "Bearer " prefix
				apiKey = strings.TrimPrefix(apiKey, "Bearer ")

				// Validate the API key
				if err := validate(apiKey, config); err != nil {
					utils.WriteJSONError(w, err)
					return
				}

				next.ServeHTTP(w, r)
			})
		},
	}
}
