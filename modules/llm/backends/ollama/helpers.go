package ollama

import "github.com/YourTechBud/inferix/modules/llm/types"

func convertToOllamaRequest(req types.InferenceRequest, opts types.InferenceOptions) OllamaRequest {
	ollamaOptions := make(map[string]any)
	if opts.NumCtx != nil {
		ollamaOptions["num_ctx"] = *opts.NumCtx
	}
	if opts.Temperature != nil {
		ollamaOptions["temperature"] = *opts.Temperature
	}
	if opts.TopK != nil {
		ollamaOptions["top_k"] = *opts.TopK
	}
	if opts.TopP != nil {
		ollamaOptions["top_p"] = *opts.TopP
	}

	ollamaMessages := make([]OllamaMessage, len(req.Messages))
	for i, msg := range req.Messages {
		ollamaMessages[i] = OllamaMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	return OllamaRequest{
		Model:    req.Model,
		Messages: ollamaMessages,
		Options:  ollamaOptions,
		Stream:   false,
	}
}
