package llm

import (
	"github.com/go-chi/chi/v5"

	"github.com/YourTechBud/inferix/modules/llm/apis/openai"
	"github.com/YourTechBud/inferix/modules/llm/apis/tei"
	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/models"
	"github.com/YourTechBud/inferix/utils"
)

// Routes defines the routes for the LLM module
func (llm *LLM) Routes() chi.Router {
	return llm.routes
}

// Middlewares defines the middlewares for the LLM module
func (llm *LLM) Middlewares() []utils.HTTPMiddleware {
	// The llm module does not have any middlewares
	return nil
}

func initializeRoutes(models *models.Models, backends *backends.Backends) chi.Router {
	router := chi.NewRouter()

	// APIs for OpenAI
	router.Post("/chat/completions", openai.HandleChatCompletion(backends))
	router.Post("/embeddings", openai.HandleCreateEmbeddings(backends))
	router.Get("/models", openai.HandleGetModels(models))

	// APIs for TEI
	router.Post("/embed", tei.HandleEmbed(backends))

	// Return the router
	return router
}
