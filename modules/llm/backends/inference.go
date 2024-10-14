package backends

import (
	"context"
	"errors"
	"strings"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// RunInference is a function that runs inference on the backend
func (b *Backends) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponseSync, error) {
	// Get the model first
	modelConfig, err := b.models.GetModel(req.Model)
	if err != nil {
		return types.InferenceResponseSync{}, err
	}

	// Set the target name of the model and update its options with the default ones.
	req.Model = modelConfig.GetTargetName()
	modelConfig.MergeOptions(&opts)

	// TODO: Inject function calling prompt if the tools are provided and backend doesn't support it

	// Get the right backend
	backend, err := b.getBackend(modelConfig.Driver) // TODO: Get the backend based on the model
	if err != nil {
		return types.InferenceResponseSync{}, err
	}

	// Run inference in a loop
	// TODO: Make the retry count configurable
	for i := 0; i < 3; i++ {
		resp, err := backend.RunInference(ctx, req, opts)
		if err != nil {
			return types.InferenceResponseSync{}, err
		}

		// Remove all whitespaces from the response
		resp.Response.Content = strings.TrimSpace(resp.Response.Content)

		// Check if the response is not too short
		if len(resp.Response.Content) < 5 {
			continue
		}

		// Generate an id if it doesn't already exist
		if resp.ID == "" {
			resp.ID = "inferix" // TODO: Generate a unique id
		}

		// Inject a default finish reason if it doesn't exist
		if resp.Response.FinishReason == "" {
			resp.Response.FinishReason = types.FinishReason_Stop
		}

		// TODO: Check for function call in the response if backend did not support it natively

		return resp, nil
	}

	return types.InferenceResponseSync{}, errors.New("unable to generate a response")
}
