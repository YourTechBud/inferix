package llm

import (
	"github.com/YourTechBud/inferix/modules/llm/apis/openai"
	"github.com/go-chi/chi/v5"
)

// Routes defines the routes for the LLM module
func (llm *LLM) Routes() chi.Router {
	router := chi.NewRouter()
	router.Post("/chat/completions", openai.HandleChatCompletion(llm.Backends))
	return router
}
