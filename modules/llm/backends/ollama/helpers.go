package ollama

import (
	"encoding/json"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

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

	// Add the tools to the ollama request
	var ollamaTools []OllamaFunctionCallRequest = nil
	if len(req.Tools) > 0 {
		ollamaTools = make([]OllamaFunctionCallRequest, len(req.Tools))
		for i, tool := range req.Tools {
			ollamaTools[i] = OllamaFunctionCallRequest{
				Type: "function",
				Function: FunctionCallRequest{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters: struct {
						Type       string          "json:\"type\""
						Properties json.RawMessage "json:\"properties\""
						Required   []string        "json:\"required\""
					}{
						Type:       "object",
						Properties: tool.Args,
					},
				},
			}
		}
	}

	return OllamaRequest{
		Model:    req.Model,
		Messages: ollamaMessages,
		Options:  ollamaOptions,
		Stream:   false,
		Tools:    ollamaTools,
	}
}
