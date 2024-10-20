package llm

import (
	"github.com/YourTechBud/inferix/modules/llm/apis/openai"
	"github.com/YourTechBud/inferix/modules/llm/apis/tei"
	"github.com/go-chi/chi/v5"
)

// Routes defines the routes for the LLM module
func (llm *LLM) Routes() chi.Router {
	router := chi.NewRouter()

	// APIs for OpenAI
	router.Post("/chat/completions", openai.HandleChatCompletion(llm.Backends))
	router.Post("/embeddings", openai.HandleCreateEmbeddings(llm.Backends))

	// APIs for TEI
	router.Post("/embed", tei.HandleEmbed(llm.Backends))

	// Return the router
	return router
}
