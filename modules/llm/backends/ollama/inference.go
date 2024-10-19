package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func (backend *Ollama) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponse, error) {
	ollamaReq := convertToOllamaRequest(req, opts)

	// Make http request
	httpResponse, err := utils.MakeHTTPRequest[OllamaResponse, utils.StandardError](ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", backend.BaseURL), ollamaReq)
	if err != nil {
		err = utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_call_error")
		return types.InferenceResponse{}, err
	}

	// Check if we have an error
	if httpResponse.HasError() {
		// Set the status of the error to the status of the response
		httpResponse.Error.Status = httpResponse.Status

		// Return the error
		return types.InferenceResponse{}, *httpResponse.Error
	}

	// Convert the created at time
	createdAt, err := time.Parse(time.RFC3339, httpResponse.Data.CreatedAt)
	if err != nil {
		err = utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_invalid_response")
		return types.InferenceResponse{}, err
	}

	// Get the tool calls
	var fnCall *types.FunctionCall
	if len(httpResponse.Data.Message.ToolCalls) > 0 {
		fnCall = &types.FunctionCall{
			Name:       httpResponse.Data.Message.ToolCalls[0].Function.Name,
			Parameters: httpResponse.Data.Message.ToolCalls[0].Function.Arguments,
		}
	}

	ollamaResp := httpResponse.Data
	return types.InferenceResponse{
		Model: ollamaResp.Model,
		Response: types.InferenceResponseMessage{
			Content:      ollamaResp.Message.Content,
			FinishReason: types.FinishReason_Stop,
			FnCall:       fnCall,
		},
		CreatedAt: createdAt,
		Stats: &types.InferenceStats{
			TotalDuration:      *ollamaResp.TotalDuration,
			LoadDuration:       *ollamaResp.LoadDuration,
			PromptEvalCount:    *ollamaResp.PromptEvalCount,
			PromptEvalDuration: *ollamaResp.PromptEvalDuration,
			EvalCount:          *ollamaResp.EvalCount,
			EvalDuration:       *ollamaResp.EvalDuration,
		},
	}, nil
}

func (backend *Ollama) RunStreamingInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) types.StreamingInferenceResponse {
	return func(yield func(element types.InferenceStreamingResponse) bool) {
		ollamaRequest := convertToOllamaRequest(req, opts)
		ollamaRequest.Stream = true

		// Make http request
		stream := utils.MakeHTTPStream(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", backend.BaseURL), ollamaRequest)
		for element := range stream {
			if element.Error != nil {
				yield(types.InferenceStreamingResponse{Err: element.Error})
				return
			}

			// Check if we got an invalid response
			if element.Status >= 300 && element.Status < 200 {
				data := new(utils.StandardError)
				if err := json.Unmarshal([]byte(element.Message), data); err != nil {
					yield(types.InferenceStreamingResponse{Err: utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_invalid_response")})
					return
				}

				data.Status = element.Status
				yield(types.InferenceStreamingResponse{Err: data})
				return
			}

			// Convert the response to an ollama response
			ollamaResp := new(OllamaResponse)
			if err := json.Unmarshal([]byte(element.Message), ollamaResp); err != nil {
				yield(types.InferenceStreamingResponse{Err: utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_invalid_response")})
				return
			}

			// Convert the created at time
			createdAt, err := time.Parse(time.RFC3339, ollamaResp.CreatedAt)
			if err != nil {
				yield(types.InferenceStreamingResponse{Err: utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_invalid_response")})
				return
			}

			// Get the message
			message := types.InferenceResponseMessage{}
			if !ollamaResp.Done {
				message.Content = ollamaResp.Message.Content
			}

			// Get the stats
			var stats *types.InferenceStats = nil
			if ollamaResp.Done {
				message.FinishReason = types.FinishReason_Stop
				stats = &types.InferenceStats{
					TotalDuration:      *ollamaResp.TotalDuration,
					LoadDuration:       *ollamaResp.LoadDuration,
					PromptEvalCount:    *ollamaResp.PromptEvalCount,
					PromptEvalDuration: *ollamaResp.PromptEvalDuration,
					EvalCount:          *ollamaResp.EvalCount,
					EvalDuration:       *ollamaResp.EvalDuration,
				}

				// Need to send one additional message with the stats when done
				if !yield(types.InferenceStreamingResponse{
					Data: types.InferenceResponse{
						Model:     ollamaResp.Model,
						Done:      ollamaResp.Done,
						CreatedAt: createdAt,
						Response:  message,
						Stats:     stats,
					},
				}) {
					return
				}
			}

			if !yield(types.InferenceStreamingResponse{
				Data: types.InferenceResponse{
					Model:     ollamaResp.Model,
					Done:      ollamaResp.Done,
					CreatedAt: createdAt,
					Response:  message,
					Stats:     stats,
				},
			}) {
				return
			}
		}
	}
}
