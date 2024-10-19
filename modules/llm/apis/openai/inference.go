package openai

import (
	"encoding/json"
	"fmt"
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

		// Don't allow tool use for streaming responses
		if req.Stream && (len(req.Tools) > 0 || len(req.Functions) > 0) {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "Tools and functions are not allowed for streaming responses", "invalid_request"))
			return
		}

		// Prepare the tools
		toolSelection := "none"
		var tools []types.Tool

		if len(req.Tools) > 0 {
			toolSelection = "tool"
			tools = make([]types.Tool, len(req.Tools))
			for i, tool := range req.Tools {
				tools[i] = types.Tool{
					Name:        tool.Function.Name,
					Description: tool.Function.Description,
					Args:        tool.Function.Parameters,
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
		if req.Stream {
			stream := backends.RunStreamingInference(r.Context(), inferenceRequest, opts)

			var once, isDone bool

			for element := range stream {
				// Handle the error first
				if element.Err != nil {
					utils.WriteJSONError(w, element.Err)
					return
				}

				// Set the response headers
				if !once {
					w.Header().Set("Content-Type", "text/event-stream")
					w.Header().Set("Cache-Control", "no-cache")
					w.Header().Set("Connection", "keep-alive")

					once = true
				}

				// Prepare the response
				delta := types.ChatCompletionStreamResponseDelta{}
				delta.Content = element.Data.Response.Content
				delta.Role = "assistant"

				var usage *types.CompletionUsage = nil
				if element.Data.Stats != nil {
					usage = &types.CompletionUsage{
						CompletionTokens: element.Data.Stats.EvalCount,
						PromptTokens:     element.Data.Stats.PromptEvalCount,
						TotalTokens:      element.Data.Stats.EvalCount + element.Data.Stats.PromptEvalCount,
					}
				}

				// Prepare the choices
				choices := make([]types.StreamingResponseChoice, 0, 1)
				if !isDone {
					choices = append(choices, types.StreamingResponseChoice{
						Index:        0,
						Delta:        delta,
						FinishReason: element.Data.Response.FinishReason,
					})

					isDone = element.Data.Done
				}

				response := types.CreateChatCompletionStreamResponse{
					ID:      element.Data.ID,
					Model:   element.Data.Model,
					Object:  "chat.completion",
					Created: element.Data.CreatedAt.Unix(),
					Choices: choices,
					Usage:   usage,
				}

				data, _ := json.Marshal(response)
				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			}

			// Send the done message
			fmt.Fprint(w, "data: [DONE]\n\n")
			w.(http.Flusher).Flush()
			return
		}

		res, err := backends.RunInference(r.Context(), inferenceRequest, opts)
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
				prepareResponseChoices(res.Response, toolSelection),
			},
		}

		// Write the response
		utils.WriteJSON(w, response)
	}
}
