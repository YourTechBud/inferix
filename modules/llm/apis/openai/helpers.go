package openai

import "github.com/YourTechBud/inferix/modules/llm/types"

func prepareResponseChoices(response types.InferenceResponseMessage) types.ResponseChoice {
	return types.ResponseChoice{
		Index:        0,
		FinishReason: response.FinishReason,
		Message: types.ChatCompletionResponseMessage{
			Role:    "assistant",
			Content: response.Content,
			// TODO: Add support for tool and function calls
		},
	}
}
