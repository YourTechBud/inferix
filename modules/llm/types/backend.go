package types

import (
	"context"
	"io"
)

// Backend is an interface for all backends
type Backend interface {
	RunInference(ctx context.Context, req InferenceRequest, opts InferenceOptions) (InferenceResponse, error)
	RunStreamingInference(ctx context.Context, req InferenceRequest, opts InferenceOptions) StreamingInferenceResponse

	CreateEmbeddings(ctx context.Context, req EmbeddingRequest) (EmbeddingResponse, error)

	TextToSpeech(ctx context.Context, req TextToSpeechRequest) (io.Reader, error)

	GetModels(ctx context.Context) ([]ModelObject, error)

	RunFnInjection() bool
	EnableEmbeddingsAPI() bool
}

type StreamingInferenceResponse func(yield func(element InferenceStreamingResponse) bool)
