package openai

import (
	"github.com/YourTechBud/inferix/modules/llm/types"
)

func convertToOpenAIRequest(req types.InferenceRequest, opts types.InferenceOptions) types.CreateChatCompletionRequest {

	// Convert messages
	messages := make([]types.ChatCompletionRequestMessage, len(req.Messages))
	for i, message := range req.Messages {
		messages[i] = types.ChatCompletionRequestMessage{
			Role:    message.Role,
			Content: message.Content,

			// TODO: Add tools for backends which support it
		}
	}
	return types.CreateChatCompletionRequest{
		Model:     req.Model,
		Messages:  messages,
		MaxTokens: opts.NumCtx,
		// TODO: Add support for MaxCompletionToken based on model config

		Temperature: opts.Temperature,
		TopP:        opts.TopP,
		N:           1,
		Stream:      false,

		// TODO: Add support for tools for backends which support it
		// TODO: Add support for structure response for backends which support it
	}
}
