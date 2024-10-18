package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func (backend *OpenAI) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponseSync, error) {
	openaiReq := convertToOpenAIRequest(req, opts)

	// Make http request
	// TODO: Add support for api key
	httpResponse, err := utils.MakeHTTPRequest[types.CreateChatCompletionResponse, utils.StandardError](ctx, http.MethodPost, fmt.Sprintf("%s/chat/completions", backend.BaseURL), openaiReq)
	if err != nil {
		err = utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_call_error")
		return types.InferenceResponseSync{}, err
	}

	// Check if we have an error
	if httpResponse.HasError() {
		// Set the status of the error to the status of the response
		httpResponse.Error.Status = httpResponse.Status

		// Return the error
		return types.InferenceResponseSync{}, *httpResponse.Error
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

	return types.InferenceResponseSync{
		ID:    openaiResp.ID,
		Model: openaiReq.Model,
		Response: types.InferenceResponseMessage{
			Content:      openaiResp.Choices[0].Message.Content,
			FinishReason: finishReason,
			FnCall:       fnCall,
		},
		CreatedAt: time.Unix(openaiResp.Created, 0),
		Stats: types.InferenceStats{
			EvalCount:       uint64(openaiResp.Usage.CompletionTokens),
			PromptEvalCount: uint64(openaiResp.Usage.PromptTokens),

			// TODO: Add support for total duration, load duration, and prompt eval duration
		},
	}, nil
}
