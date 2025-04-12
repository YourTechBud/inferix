package openai

import (
	"context"
	"io"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// TextToSpeech converts text to speech using the OpenAI API
func (b *OpenAI) TextToSpeech(ctx context.Context, req types.TextToSpeechRequest) (io.Reader, error) {
	// Make the request to the OpenAI API
	reader, err := utils.MakeHTTPBinaryStream(ctx, http.MethodPost, b.BaseURL+"/audio/speech", req)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
