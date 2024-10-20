package tei

import (
	"encoding/json"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func HandleEmbed(backends *backends.Backends) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create embeddings
		embeddingsRequest := types.EmbeddingRequest{Model: "default", Input: req.Inputs}
		embeddingsResponse, err := backends.CreateEmbeddings(r.Context(), embeddingsRequest)
		if err != nil {
			utils.WriteJSONError(w, err)
			return

		}

		// Don't forget to convert the response to the right format
		res := make(EmbeddingResponse, len(embeddingsResponse.Data))
		for i, data := range embeddingsResponse.Data {
			res[i] = data.Embedding
		}

		// Write the response
		utils.WriteJSON(w, res)
	}
}
