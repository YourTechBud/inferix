package openai

import (
	"context"
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
	return types.InferenceResponseSync{
		ID:    openaiResp.ID,
		Model: openaiReq.Model,
		Response: types.InferenceResponseMessage{
			Content:      openaiResp.Choices[0].Message.Content,
			FinishReason: types.FinishReason(openaiResp.Choices[0].FinishReason),
			// TODO: Add support to return tool calls for backends which support it
		},
		CreatedAt: time.Unix(openaiResp.Created, 0),
		Stats: types.InferenceStats{
			EvalCount:       uint64(openaiResp.Usage.CompletionTokens),
			PromptEvalCount: uint64(openaiResp.Usage.PromptTokens),

			// TODO: Add support for total duration, load duration, and prompt eval duration
		},
	}, nil
}
