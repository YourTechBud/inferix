package openai

import (
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/YourTechBud/inferix/utils"
)

// HandleGetModels returns a handler that returns the list of models
func HandleGetModels(models *models.Models) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resources, err := models.GetModels(r.Context())
		if err != nil {
			utils.WriteJSONError(w, err)
			return
		}

		utils.WriteJSON(w, map[string]any{
			"object": "list",
			"data":   resources,
		})
	}
}
