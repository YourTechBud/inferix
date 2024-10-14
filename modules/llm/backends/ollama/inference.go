package ollama

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func (backend *Ollama) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponseSync, error) {
	ollamaReq := convertToOllamaRequest(req, opts)

	// Make http request
	httpResponse, err := utils.MakeHTTPRequest[OllamaResponse, utils.StandardError](ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", backend.BaseURL), ollamaReq)
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

	// Convert the created at time
	createdAt, err := time.Parse(time.RFC3339, httpResponse.Data.CreatedAt)
	if err != nil {
		err = utils.NewStandardError(http.StatusInternalServerError, err.Error(), "backend_invalid_response")
		return types.InferenceResponseSync{}, err
	}

	ollamaResp := httpResponse.Data
	return types.InferenceResponseSync{
		Model: ollamaResp.Model,
		Response: types.InferenceResponseMessage{
			Content:      ollamaResp.Message.Content,
			FinishReason: types.FinishReason_Stop,
		},
		CreatedAt: createdAt,
		Stats: types.InferenceStats{
			TotalDuration:      *ollamaResp.TotalDuration,
			LoadDuration:       *ollamaResp.LoadDuration,
			PromptEvalCount:    *ollamaResp.PromptEvalCount,
			PromptEvalDuration: *ollamaResp.PromptEvalDuration,
			EvalCount:          *ollamaResp.EvalCount,
			EvalDuration:       *ollamaResp.EvalDuration,
		},
	}, nil
}
