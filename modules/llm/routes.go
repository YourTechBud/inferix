package llm

import (
	"github.com/YourTechBud/inferix/modules/llm/apis/openai"
	"github.com/YourTechBud/inferix/modules/llm/apis/tei"
	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/go-chi/chi/v5"
)

// Routes defines the routes for the LLM module
func (llm *LLM) Routes() chi.Router {
	return llm.routes
}

func initializeRoutes(backends *backends.Backends) chi.Router {
	router := chi.NewRouter()

	// APIs for OpenAI
	router.Post("/chat/completions", openai.HandleChatCompletion(backends))
	router.Post("/embeddings", openai.HandleCreateEmbeddings(backends))

	// APIs for TEI
	router.Post("/embed", tei.HandleEmbed(backends))

	// Return the router
	return router
}
