package backends

import (
	"context"
	"io"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

// TextToSpeech generates audio from text
func (b *Backends) TextToSpeech(ctx context.Context, req types.TextToSpeechRequest) (io.Reader, error) {
	// Find a backend based on the models
	model, err := b.models.GetModel(req.Model)
	if err != nil {
		return nil, err
	}

	// Get the backend
	backend, err := b.getBackend(model.Backend)
	if err != nil {
		return nil, err
	}

	// Call the backend's TextToSpeech method
	return backend.TextToSpeech(ctx, req)
}
