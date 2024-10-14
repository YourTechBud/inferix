package types

import (
	"context"
)

// Backend is an interface for all backends
type Backend interface {
	RunInference(ctx context.Context, req InferenceRequest, opts InferenceOptions) (InferenceResponseSync, error)
}
