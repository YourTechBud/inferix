package tei

import (
	"context"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// RunInference is not implemented for the TEI backend
func (backend *TEI) RunInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) (types.InferenceResponse, error) {
	return types.InferenceResponse{}, utils.NewStandardError(
		http.StatusNotImplemented,
		"RunInference is not implemented for the TEI backend",
		"not_implemented",
	)
}

// RunStreamingInference is not implemented for the TEI backend“
func (backend *TEI) RunStreamingInference(ctx context.Context, req types.InferenceRequest, opts types.InferenceOptions) types.StreamingInferenceResponse {
	return func(yield func(element types.InferenceStreamingResponse) bool) {
		yield(types.InferenceStreamingResponse{
			Err: utils.NewStandardError(
				http.StatusNotImplemented,
				"RunStreamingInference is not implemented for the TEI backend",
				"not_implemented",
			),
		})
	}
}
