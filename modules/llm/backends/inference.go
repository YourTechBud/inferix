package backends

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// RunInference is a function that runs inference on the backend
func (b *Backends) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponse, error) {
	// Get the model and backend config
	modelConfig, backend, err := b.getModelAndBackend(&req, &opts)
	if err != nil {
		return types.InferenceResponse{}, err
	}

	// Inject function calling prompt if the tools are provided and backend doesn't support it
	if backend.RunFnInjection() && len(req.Tools) > 0 {
		req.Messages = injectFnCall(req.Messages, req.Tools)

		// Remove the tools from the request
		req.Tools = nil
	}

	// Run inference in a loop
	// TODO: Make the retry count configurable
	for i := 0; i < 3; i++ {
		resp, err := backend.RunInference(ctx, req, opts)
		if err != nil {
			log.Default().Printf("unable to run inference: %v", err)
			return types.InferenceResponse{}, err
		}

		// Process the response
		resp = processInferenceResponse(resp, modelConfig)

		// Non-streaming processing
		// Remove all whitespaces from the response
		resp.Response.Content = strings.TrimSpace(resp.Response.Content)

		// Check if the response is too short
		if len(resp.Response.Content) < 5 && resp.Response.FnCall == nil {
			continue
		}

		// Check for function call in the response if backend did not support it natively
		if backend.RunFnInjection() {
			if strings.Contains(resp.Response.Content, "FUNC_CALL") {
				// Sanitize the JSON text
				content := utils.SanitizeJSONText(resp.Response.Content)

				// Check if response is valid JSON
				var fnCall types.FunctionCall
				if err := json.Unmarshal([]byte(content), &fnCall); err != nil {
					log.Default().Printf("unable to unmarshal function call: %v", err)
					continue
				}

				// Set the function call parameters in the response
				resp.Response.Content = fmt.Sprintf("Execute function %s with arguments: %s", fnCall.Name, string(fnCall.Parameters))
				resp.Response.FnCall = &fnCall
				resp.Response.FinishReason = types.FinishReason_FunctionCall

				// Return the response
				return resp, nil
			}
		}

		return resp, nil
	}

	return types.InferenceResponse{}, errors.New("unable to generate a response")
}

func (b *Backends) RunStreamingInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) types.StreamingInferenceResponse {
	return types.StreamingInferenceResponse(func(yield func(element types.InferenceStreamingResponse) bool) {
		// Get the model and backend config
		modelConfig, backend, err := b.getModelAndBackend(&req, &opts)
		if err != nil {
			yield(types.InferenceStreamingResponse{Err: err})
			return
		}

		// Run the inference
		stream := backend.RunStreamingInference(ctx, req, opts)
		for element := range stream {
			// Check if we have an error
			if element.Err != nil {
				yield(types.InferenceStreamingResponse{Err: element.Err})
				return
			}

			// Yield the response post processing
			resp := processInferenceResponse(element.Data, modelConfig)
			if !yield(types.InferenceStreamingResponse{Data: resp}) {
				return
			}
		}
	})
}
