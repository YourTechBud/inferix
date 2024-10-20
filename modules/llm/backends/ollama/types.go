package ollama

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// OllamaRequest represents the structure for the request.
type OllamaRequest struct {
	Model    string                      `json:"model"`
	Messages []OllamaMessage             `json:"messages"`
	Tools    []OllamaFunctionCallRequest `json:"tools,omitempty"`
	// TODO: Add support for tools

	// Optional fields
	Stream  bool                `json:"stream"`
	Options map[string]any      `json:"options"` // Raw JSON value
	Format  *types.OutputFormat `json:"format,omitempty"`
}

type OllamaMessage struct {
	Role      string                       `json:"role"`
	Content   string                       `json:"content,omitempty"`
	ToolCalls []OllamaFunctionCallResponse `json:"tool_calls,omitempty"`

	// TODO: Add support for images
}

type OllamaFunctionCallRequest struct {
	Type     string              `json:"type"`
	Function FunctionCallRequest `json:"function"`
}

type FunctionCallRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string          `json:"type"`
		Properties json.RawMessage `json:"properties"`
		Required   []string        `json:"required"`
	} `json:"parameters"`
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

type OllamaFunctionCallResponse struct {
	Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	} `json:"function"`
}

// Types for creating embeddings
type OllamaEmbeddingsRequest struct {
	types.EmbeddingRequest

	// Optional fields
	Truncate bool           `json:"truncate,omitempty"`
	Options  map[string]any `json:"options,omitempty"` // Raw JSON value
}

type OllamaEmbeddingsResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float32 `json:"embeddings"`
	TotalDuration   uint64      `json:"total_duration"`
	LoadDuration    uint64      `json:"load_duration"`
	PromptEvalCount uint64      `json:"prompt_eval_count"`
}
