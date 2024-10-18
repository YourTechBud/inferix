package openai

import "github.com/YourTechBud/inferix/modules/llm/types"

func prepareResponseChoices(response types.InferenceResponseMessage, toolSelection string) types.ResponseChoice {
	// Declare variables for the tools and function calls
	var fnCall *types.ChatCompletionFunctionCall = nil
	var toolCalls []types.ChatCompletionMessageToolCall

	// Check if there are any function calls being made
	if response.FnCall != nil {
		switch toolSelection {
		case "tool":
			toolCalls = []types.ChatCompletionMessageToolCall{
				types.ChatCompletionMessageToolCall{
					ID:       response.FnCall.Name,
					ToolType: types.Function,
					Function: types.ChatCompletionFunctionCall{
						Name:      response.FnCall.Name,
						Arguments: string(response.FnCall.Parameters),
					},
				},
			}
			response.FinishReason = types.FinishReason_ToolCalls

		case "function":
			fnCall = &types.ChatCompletionFunctionCall{
				Name:      response.FnCall.Name,
				Arguments: string(response.FnCall.Parameters),
			}
			response.FinishReason = types.FinishReason_FunctionCall
		}
	}
	return types.ResponseChoice{
		Index:        0,
		FinishReason: response.FinishReason,
		Message: types.ChatCompletionResponseMessage{
			Role:         "assistant",
			Content:      response.Content,
			FunctionCall: fnCall,
			ToolCalls:    toolCalls,
		},
	}
}
