package types

import (
	"encoding/json"
	"time"
)

type InferenceRequest struct {
	Model        string             `json:"model"`
	Messages     []InferenceMessage `json:"messages"`
	Tools        []Tool             `json:"tools,omitempty"`
	OutputFormat *OutputFormat      `json:"format,omitempty"`
}

func NewInferenceRequest(model string, messages []InferenceMessage, tools []Tool) *InferenceRequest {
	return &InferenceRequest{
		Model:    model,
		Messages: messages,
		Tools:    tools,
	}
}

type OutputFormat string

const (
	OutputFormat_JSON OutputFormat = "json"
)

type InferenceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
}

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Args        json.RawMessage `json:"args"`
	ToolType    ToolType        `json:"type"`
}

type ToolType string

const (
	Function ToolType = "function"
)

type InferenceOptions struct {
	TopP          *float64        `json:"top_p,omitempty"`
	TopK          *int32          `json:"top_k,omitempty"`
	NumCtx        *int32          `json:"num_ctx,omitempty"`
	Temperature   *float64        `json:"temperature,omitempty"`
	DriverOptions json.RawMessage `json:"driver_options"`
}

func NewInferenceOptions(topP *float64, topK *int32, numCtx *int32, temperature *float64) InferenceOptions {
	return InferenceOptions{
		TopP:          topP,
		TopK:          topK,
		NumCtx:        numCtx,
		Temperature:   temperature,
		DriverOptions: json.RawMessage(`{}`),
	}
}

func DefaultInferenceOptions() InferenceOptions {
	topP := 0.9
	topK := int32(40)
	numCtx := int32(4096)
	temperature := 0.2

	return InferenceOptions{
		TopP:          &topP,
		TopK:          &topK,
		NumCtx:        &numCtx,
		Temperature:   &temperature,
		DriverOptions: json.RawMessage(`{}`),
	}
}

type InferenceResponseSync struct {
	ID        string                   `json:"id"`
	Model     string                   `json:"model"`
	CreatedAt time.Time                `json:"created_at"`
	Response  InferenceResponseMessage `json:"response"`
	Stats     InferenceStats           `json:"stats"`
}

type InferenceResponseStream struct {
	Model        string                   `json:"model"`
	CreatedAt    time.Time                `json:"created_at"`
	Response     InferenceResponseMessage `json:"response"`
	FinishReason *FinishReason            `json:"finish_reason,omitempty"`
	Stats        *InferenceStats          `json:"stats,omitempty"`
}

type InferenceResponseMessage struct {
	Content      string        `json:"content"`
	FnCall       *FunctionCall `json:"fn_call,omitempty"`
	FinishReason FinishReason  `json:"finish_reason,omitempty"`
}

type FinishReason string

const (
	FinishReason_Stop          FinishReason = "stop"
	FinishReason_Length        FinishReason = "length"
	FinishReason_ToolCalls     FinishReason = "tool_calls"
	FinishReason_ContentFilter FinishReason = "content_filter"
	FinishReason_FunctionCall  FinishReason = "function_call"
)

type FunctionCall struct {
	Name       string          `json:"name"`
	Parameters json.RawMessage `json:"parameters"`
}

type InferenceStats struct {
	TotalDuration      uint64 `json:"total_duration,omitempty"`
	LoadDuration       uint64 `json:"load_duration,omitempty"`
	PromptEvalCount    uint64 `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration uint64 `json:"prompt_eval_duration,omitempty"`
	EvalCount          uint64 `json:"eval_count,omitempty"`
	EvalDuration       uint64 `json:"eval_duration,omitempty"`
}
