package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func (backend *OpenAI) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponse, error) {
	openaiReq := convertToOpenAIRequest(req, opts)

	// Make http request
	// TODO: Add support for api key
	httpResponse, err := utils.MakeHTTPRequest[types.CreateChatCompletionResponse, utils.StandardError](ctx, http.MethodPost, fmt.Sprintf("%s/chat/completions", backend.BaseURL), openaiReq)
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

	openaiResp := httpResponse.Data

	// Set the finish reason
	finishReason := types.FinishReason(openaiResp.Choices[0].FinishReason)

	// Get the tools
	var fnCall *types.FunctionCall = nil
	if len(openaiResp.Choices[0].Message.ToolCalls) > 0 {
		fnCall = &types.FunctionCall{
			Name:       openaiResp.Choices[0].Message.ToolCalls[0].Function.Name,
			Parameters: json.RawMessage(openaiResp.Choices[0].Message.ToolCalls[0].Function.Arguments),
		}
		finishReason = types.FinishReason_ToolCalls
	}

	return types.InferenceResponse{
		ID:    openaiResp.ID,
		Model: openaiReq.Model,
		Response: types.InferenceResponseMessage{
			Content:      openaiResp.Choices[0].Message.Content,
			FinishReason: finishReason,
			FnCall:       fnCall,
		},
		CreatedAt: time.Unix(openaiResp.Created, 0),
		Stats: &types.InferenceStats{
			EvalCount:       uint64(openaiResp.Usage.CompletionTokens),
			PromptEvalCount: uint64(openaiResp.Usage.PromptTokens),

			// TODO: Add support for total duration, load duration, and prompt eval duration
		},
	}, nil
}

func (backend *OpenAI) RunStreamingInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) types.StreamingInferenceResponse {
	return func(yield func(element types.InferenceStreamingResponse) bool) {
		openaiReq := convertToOpenAIRequest(req, opts)
		openaiReq.Stream = true
		openaiReq.StreamOptions = &types.ChatCompletionStreamOptions{IncludeUsage: true}

		// Make http request
		stream := utils.MakeHTTPStream(ctx, http.MethodPost, fmt.Sprintf("%s/chat/completions", backend.BaseURL), openaiReq)
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

			// Strip the data: prefix from the SSE message
			data := strings.TrimPrefix(element.Message, "data: ")

			// Skip if it's the last message
			if strings.Contains(data, "[DONE]") {
				continue
			}

			// Parse the response
			openaiResp := new(types.CreateChatCompletionStreamResponse)
			if err := json.Unmarshal([]byte(data), openaiResp); err != nil {
				yield(types.InferenceStreamingResponse{Err: utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_invalid_response")})
				return
			}

			// Get the message
			message := types.InferenceResponseMessage{}
			if len(openaiResp.Choices) > 0 {
				message.Content = openaiResp.Choices[0].Delta.Content
				message.FinishReason = types.FinishReason(openaiResp.Choices[0].FinishReason)
			}

			// Check if the response is done
			isDone := message.FinishReason == types.FinishReason_Stop || len(openaiResp.Choices) == 0

			// Get the stats
			var stats *types.InferenceStats = nil
			if openaiResp.Usage != nil {
				stats = &types.InferenceStats{
					EvalCount:       uint64(openaiResp.Usage.CompletionTokens),
					PromptEvalCount: uint64(openaiResp.Usage.PromptTokens),
				}
			}

			yield(types.InferenceStreamingResponse{
				Data: types.InferenceResponse{
					ID:        openaiResp.ID,
					Model:     openaiResp.Model,
					CreatedAt: time.Unix(openaiResp.Created, 0),
					Response:  message,
					Stats:     stats,
					Done:      isDone,
				},
			})
		}
	}
}
