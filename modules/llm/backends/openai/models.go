package openai

import (
	"context"
	"fmt"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

type modelsResponse struct {
	Data []types.ModelObject `json:"data"`
}

// GetModels retrieves the list of available models from the OpenAI API
func (o *OpenAI) GetModels(ctx context.Context) ([]types.ModelObject, error) {
	// Make HTTP request to OpenAI API
	resp, err := utils.MakeHTTPRequest[modelsResponse, any](ctx, "GET", fmt.Sprintf("%s/models", o.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models from OpenAI: %w", err)
	}

	if resp.HasError() {
		return nil, fmt.Errorf("failed to fetch models from OpenAI: %v", resp.Error)
	}

	return resp.Data.Data, nil
}
