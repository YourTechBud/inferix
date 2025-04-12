package types

// TextToSpeechRequest is the request body for the text-to-speech API
type TextToSpeechRequest struct {
	Input          string  `json:"input"`
	Model          string  `json:"model"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed"`
}

// DefaultTextToSpeechResponseFormat is the default response format for text-to-speech
const DefaultTextToSpeechResponseFormat = "wav"
