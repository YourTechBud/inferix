package ollama

import (
	"context"
	"io"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// TextToSpeech is not implemented for the Ollama backend
func (b *Ollama) TextToSpeech(ctx context.Context, req types.TextToSpeechRequest) (io.Reader, error) {
	return nil, utils.NewStandardError(http.StatusNotImplemented, "text to speech is not implemented for the Ollama backend", "not_implemented")
}
