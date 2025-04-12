package openai

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/YourTechBud/inferix/modules/llm/backends"
	"github.com/YourTechBud/inferix/modules/llm/types"
	"github.com/YourTechBud/inferix/utils"
)

// HandleTextToSpeech handles the text-to-speech API endpoint
func HandleTextToSpeech(backends *backends.Backends) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req types.TextToSpeechRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, err.Error(), "invalid_request"))
			return
		}

		// Set default values if not provided
		if req.Model == "" {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "model is required", "invalid_request"))
			return
		}
		if req.Voice == "" {
			utils.WriteJSONError(w, utils.NewStandardError(http.StatusBadRequest, "voice is required", "invalid_request"))
			return
		}
		if req.ResponseFormat == "" {
			req.ResponseFormat = types.DefaultTextToSpeechResponseFormat
		}
		if req.Speed == 0 {
			req.Speed = 1.0
		}

		// Get backend context from the request
		ctx := r.Context()

		// Generate speech from the backend
		audioReader, err := backends.TextToSpeech(ctx, req)
		if err != nil {
			utils.WriteJSONError(w, err)
			return
		}

		// Set content type based on response format
		switch req.ResponseFormat {
		case "mp3":
			w.Header().Set("Content-Type", "audio/mpeg")
		case "opus":
			w.Header().Set("Content-Type", "audio/opus")
		case "aac":
			w.Header().Set("Content-Type", "audio/aac")
		case "flac":
			w.Header().Set("Content-Type", "audio/flac")
		default: // wav
			w.Header().Set("Content-Type", "audio/wav")
		}

		// Write the audio data directly to the response
		io.Copy(w, audioReader)
	}
}
