package openai

import (
	"encoding/json"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// HandleCreateEmbeddings handles the request to create embeddings
func HandleCreateEmbeddings(backends *backends.Backends) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req types.EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create embeddings
		response, err := backends.CreateEmbeddings(r.Context(), req)
		if err != nil {
			utils.WriteJSONError(w, err)
			return
		}

		// Write the response
		utils.WriteJSON(w, response)
	}
}
