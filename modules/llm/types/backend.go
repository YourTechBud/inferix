package types

import (
	"context"
)

// Backend is an interface for all backends
type Backend interface {
	RunInference(ctx context.Context, req InferenceRequest, opts InferenceOptions) (InferenceResponse, error)
	RunStreamingInference(ctx context.Context, req InferenceRequest, opts InferenceOptions) StreamingInferenceResponse

	CreateEmbeddings(ctx context.Context, req EmbeddingRequest) (EmbeddingResponse, error)

	RunFnInjection() bool
}

type StreamingInferenceResponse func(yield func(element InferenceStreamingResponse) bool)
