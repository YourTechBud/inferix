package tei

import (
	"context"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

func (backend *TEI) CreateEmbeddings(ctx context.Context, req types.EmbeddingRequest) (types.EmbeddingResponse, error) {
	// Make the request
	teiResponse, err := utils.MakeHTTPRequest[types.EmbeddingResponse, utils.StandardError](ctx, http.MethodPost, backend.BaseURL+"/embeddings", req)
	if err != nil {
		return types.EmbeddingResponse{}, err
	}

	// Check if we have an error
	if teiResponse.HasError() {
		// Set the status of the error to the status of the response
		teiResponse.Error.Status = teiResponse.Status

		// Return the error
		return types.EmbeddingResponse{}, teiResponse.Error
	}

	return *teiResponse.Data, nil
}
