package openai

import (
	"encoding/json"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func HandleChatCompletion(backends *backends.Backends) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req types.CreateChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Prepare the request
		messages := make([]types.InferenceMessage, len(req.Messages))
		for i, message := range req.Messages {
			messages[i] = types.InferenceMessage{
				Role:    message.Role,
				Content: message.Content,
			}
		}

		// Prepare the tools
		toolSelection := "none"
		var tools []types.Tool

		if len(req.Tools) > 0 {
			toolSelection = "tool"
			tools = make([]types.Tool, len(req.Tools))
			for i, tool := range req.Tools {
				data, _ := json.Marshal(tool.Function.Parameters)
				tools[i] = types.Tool{
					Name:        tool.Function.Name,
					Description: tool.Function.Description,
					Args:        data,
					ToolType:    types.Function,
				}
			}
		}

		if toolSelection == "none" && len(req.Functions) > 0 {
			toolSelection = "function"
			tools = make([]types.Tool, len(req.Functions))
			for i, function := range req.Functions {
				data, _ := json.Marshal(function.Parameters)
				tools[i] = types.Tool{
					Name:        function.Name,
					Description: function.Description,
					Args:        data,
					ToolType:    types.Function,
				}
			}
		}

		inferenceRequest := types.NewInferenceRequest(req.Model, messages, tools)
		opts := types.NewInferenceOptions(req.TopP, nil, req.MaxTokens, req.Temperature)

		// Run the inference
		res, err := backends.RunInference(r.Context(), *inferenceRequest, opts)
		if err != nil {
			utils.WriteJSONError(w, err)
			return
		}

		// Prepare the response
		response := types.CreateChatCompletionResponse{
			ID:      res.ID,
			Model:   res.Model,
			Created: res.CreatedAt.Unix(),
			Object:  "chat.completion",
			Usage: types.CompletionUsage{
				CompletionTokens: res.Stats.EvalCount,
				PromptTokens:     res.Stats.PromptEvalCount,
				TotalTokens:      res.Stats.EvalCount + res.Stats.PromptEvalCount,
			},
			Choices: []types.ResponseChoice{
				prepareResponseChoices(res.Response),
			},
		}

		// Write the response
		utils.WriteJSON(w, response)
	}
}
