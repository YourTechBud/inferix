package tei

import (
	"context"
	"fmt"
	"time"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

type infoResponse struct {
	ModelID        string      `json:"model_id"`
	ModelType      interface{} `json:"model_type"`
	ModelDtype     string      `json:"model_dtype"`
	ModelSha       string      `json:"model_sha"`
	MaxInputLength int         `json:"max_input_length"`
	MaxBatchTokens int         `json:"max_batch_tokens"`
	AutoTruncate   bool        `json:"auto_truncate"`
	Version        string      `json:"version"`
	DockerLabel    string      `json:"docker_label"`
	Sha            string      `json:"sha"`
}

// GetModels retrieves the model information from the TEI API
func (t *TEI) GetModels(ctx context.Context) ([]types.ModelObject, error) {
	// Make HTTP request to TEI API
	resp, err := utils.MakeHTTPRequest[infoResponse, any](ctx, "GET", fmt.Sprintf("%s/info", t.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch model info from TEI: %w", err)
	}

	if resp.HasError() {
		return nil, fmt.Errorf("failed to fetch model info from TEI: %v", resp.Error)
	}

	// Return a single model with the model ID from the info response
	return []types.ModelObject{
		{
			ID:      resp.Data.ModelID,
			Created: time.Now().Unix(), // Using current process time as creation time
			Object:  "model",
			OwnedBy: "tei",
		},
	}, nil
}
