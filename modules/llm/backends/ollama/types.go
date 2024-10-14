package ollama

import (
	"github.com/YourTechBud/inferix/modules/llm/types"
)

// OllamaRequest represents the structure for the request.
type OllamaRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	// TODO: Add support for tools

	// Optional fields
	Stream  bool                `json:"stream"`
	Options map[string]any      `json:"options"` // Raw JSON value
	Format  *types.OutputFormat `json:"format,omitempty"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`

	// TODO: Add support for tools and images
}

// OllamaResponse represents the structure for the response.
type OllamaResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	Done               bool          `json:"done"`
	TotalDuration      *uint64       `json:"total_duration,omitempty"`
	LoadDuration       *uint64       `json:"load_duration,omitempty"`
	PromptEvalCount    *uint64       `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration *uint64       `json:"prompt_eval_duration,omitempty"`
	EvalCount          *uint64       `json:"eval_count,omitempty"`
	EvalDuration       *uint64       `json:"eval_duration,omitempty"`
}
