package types

import (
	"encoding/json"
	"fmt"
)

// CreateChatCompletionRequest represents a request to create a chat completion.
type CreateChatCompletionRequest struct {
	Messages         []ChatCompletionRequestMessage  `json:"messages"`
	Model            string                          `json:"model"`
	FrequencyPenalty *float64                        `json:"frequency_penalty,omitempty"`
	MaxTokens        *int32                          `json:"max_tokens,omitempty"`
	N                int                             `json:"n,omitempty"`
	PresencePenalty  *float64                        `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat                 `json:"response_format,omitempty"`
	Seed             *int64                          `json:"seed,omitempty"`
	Stop             *Stop                           `json:"stop,omitempty"`
	Stream           bool                            `json:"stream,omitempty"`
	StreamOptions    *ChatCompletionStreamOptions    `json:"stream_options,omitempty"`
	Temperature      *float64                        `json:"temperature,omitempty"`
	TopP             *float64                        `json:"top_p,omitempty"`
	Tools            []ChatCompletionTool            `json:"tools,omitempty"`
	ToolChoice       *ChatCompletionToolChoiceOption `json:"tool_choice,omitempty"`
	FunctionCall     *FunctionCallRequest            `json:"function_call,omitempty"`
	Functions        []ChatCompletionFunctions       `json:"functions,omitempty"`
}

// ChatCompletionStreamOptions represents the options for a chat completion stream.
type ChatCompletionStreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

// ChatCompletionRequestSystemMessage represents a system message in the chat completion request.
type ChatCompletionRequestSystemMessage struct {
	Content string  `json:"content"`
	Role    string  `json:"role"`
	Name    *string `json:"name,omitempty"`
}

// ChatCompletionRequestUserMessage represents a user message in the chat completion request.
type ChatCompletionRequestUserMessage struct {
	Content string  `json:"content"`
	Role    string  `json:"role"`
	Name    *string `json:"name,omitempty"`
}

// ChatCompletionMessageToolCall represents a tool call message in the chat completion.
type ChatCompletionMessageToolCall struct {
	ID       string                     `json:"id"`
	ToolType ToolType                   `json:"type"`
	Function ChatCompletionFunctionCall `json:"function"`
}

type ChatCompletionFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatCompletionRequestToolMessage represents a tool message in the chat completion request.
type ChatCompletionRequestToolMessage struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id"`
}

// ChatCompletionRequestFunctionMessage represents a function message in the chat completion request.
type ChatCompletionRequestFunctionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
	Name    string `json:"name"`
}

// ChatCompletionRequestMessage represents a message in the chat completion request.
type ChatCompletionRequestMessage struct {
	Content      any                             `json:"content,omitempty"`
	Role         string                          `json:"role"`
	Name         *string                         `json:"name,omitempty"`
	ToolCallID   *string                         `json:"tool_call_id,omitempty"`
	ToolCalls    []ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
	FunctionCall *ChatCompletionFunctionCall     `json:"function_call,omitempty"`
}

func (message ChatCompletionRequestMessage) GetMessage() InferenceMessage {

	switch c := message.Content.(type) {
	case string:
		return InferenceMessage{
			Role:    message.Role,
			Content: c,
		}
	case []any:
		// If the content is an empty array, return an empty string
		if len(c) == 0 {
			return InferenceMessage{
				Role:    message.Role,
				Content: "",
			}
		}

		msg := InferenceMessage{Role: message.Role}

		for _, content := range c {
			m := content.(map[string]any)
			if m["type"] == "text" {
				msg.Content = m["text"].(string)
			} else if m["type"] == "image_url" {
				msg.Images = append(msg.Images, m["image_url"].(map[string]any)["url"].(string))
			}
		}

		return msg
	}

	panic(fmt.Sprintf("Unsupported content type: %T", message.Content))
}

// ChatCompletionFunctionCallOption represents an option to call a function.
type ChatCompletionFunctionCallOption struct {
	Name string `json:"name"`
}

// ChatCompletionFunctions represents a function with a name, description, and parameters.
type ChatCompletionFunctions struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Parameters  any    `json:"parameters"`
}

// FunctionObject represents a function object used in tool-based messages.
type FunctionObject struct {
	Description string          `json:"description,omitempty"`
	Name        string          `json:"name"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

// ChatCompletionTool represents a tool used in chat completion.
type ChatCompletionTool struct {
	ToolType ToolType       `json:"type"`
	Function FunctionObject `json:"function"`
}

// ChatCompletionNamedToolChoice represents a named tool choice in chat completion.
type ChatCompletionNamedToolChoice struct {
	Type     ToolType       `json:"type"`
	Function FunctionObject `json:"function"`
}

// ChatCompletionToolChoiceOption represents the tool choice options in chat completion.
type ChatCompletionToolChoiceOption string

// ResponseFormat represents the format of the response (text or JSON object).
type ResponseFormat string

const (
	TextFormat       ResponseFormat = "text"
	JsonObjectFormat ResponseFormat = "json_object"
)

// Stop represents either a string or an array of stop tokens in chat completion.
type Stop struct {
	StringVal *string   `json:"string,omitempty"`
	ArrayVal  *[]string `json:"array,omitempty"`
}

// FunctionCallRequest represents a function call request with either a string or a function call option.
type FunctionCallRequest struct {
	StringVal                        *string                           `json:"string,omitempty"`
	ChatCompletionFunctionCallOption *ChatCompletionFunctionCallOption `json:"chat_completion_function_call_option,omitempty"`
}

// CreateChatCompletionResponse represents the response to a chat completion request.
type CreateChatCompletionResponse struct {
	ID                string           `json:"id"`
	Choices           []ResponseChoice `json:"choices"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	SystemFingerprint *string          `json:"system_fingerprint,omitempty"`
	Object            string           `json:"object"`
	Usage             CompletionUsage  `json:"usage"`
}

// ResponseChoice represents a choice in the chat completion response.
type ResponseChoice struct {
	FinishReason FinishReason                  `json:"finish_reason"`
	Index        int                           `json:"index"`
	Message      ChatCompletionResponseMessage `json:"message"`
}

// ChatCompletionResponseMessage represents a message in the chat completion response.
type ChatCompletionResponseMessage struct {
	Content      string                          `json:"content,omitempty"`
	ToolCalls    []ChatCompletionMessageToolCall `json:"tool_calls,omitempty"`
	Role         string                          `json:"role"`
	FunctionCall *ChatCompletionFunctionCall     `json:"function_call,omitempty"`
}

// CompletionUsage represents usage information for the chat completion response.
type CompletionUsage struct {
	CompletionTokens uint64 `json:"completion_tokens"`
	PromptTokens     uint64 `json:"prompt_tokens"`
	TotalTokens      uint64 `json:"total_tokens"`
}

// CreateChatCompletionStreamResponse represents a streaming response for chat completion.
type CreateChatCompletionStreamResponse struct {
	ID                string                    `json:"id"`
	Choices           []StreamingResponseChoice `json:"choices"`
	Created           int64                     `json:"created"`
	Model             string                    `json:"model"`
	SystemFingerprint *string                   `json:"system_fingerprint,omitempty"`
	Object            string                    `json:"object"`
	Usage             *CompletionUsage          `json:"usage,omitempty"`
}

// StreamingResponseChoice represents a choice in the streaming response.
type StreamingResponseChoice struct {
	Delta        ChatCompletionStreamResponseDelta `json:"delta"`
	FinishReason FinishReason                      `json:"finish_reason,omitempty"`
	Index        int                               `json:"index"`
}

// ChatCompletionStreamResponseDelta represents a delta update in the chat completion stream response.
type ChatCompletionStreamResponseDelta struct {
	Role         string                                `json:"role"`
	Content      string                                `json:"content"`
	FunctionCall *ChatCompletionFunctionCall           `json:"function_call,omitempty"`
	ToolCalls    *[]ChatCompletionMessageToolCallChunk `json:"tool_calls,omitempty"`
}

// ChatCompletionMessageToolCallChunk represents a chunk of tool call message in the streaming response.
type ChatCompletionMessageToolCallChunk struct {
	Index     int                         `json:"index"`
	ID        *string                     `json:"id,omitempty"`
	TypeField *string                     `json:"type,omitempty"`
	Function  *ChatCompletionFunctionCall `json:"function,omitempty"`
}

// ModelObject describes an OpenAI model offering that can be used with the API.
type ModelObject struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}
