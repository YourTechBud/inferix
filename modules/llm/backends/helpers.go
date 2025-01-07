package backends

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YourTechBud/inferix/modules/llm/backends/ollama"
	"github.com/YourTechBud/inferix/modules/llm/backends/openai"
	"github.com/YourTechBud/inferix/modules/llm/backends/tei"
	"github.com/YourTechBud/inferix/modules/llm/config"
	"github.com/YourTechBud/inferix/modules/llm/types"
)

func (b *Backends) getModelAndBackend(req *types.InferenceRequest, opts *types.InferenceOptions) (config.ModelConfig, types.Backend, error) {
	// Get the model first
	modelConfig, err := b.models.GetModel(req.Model)
	if err != nil {
		return modelConfig, nil, err
	}

	// Set the target name of the model and update its options with the default ones.
	req.Model = modelConfig.GetTarget()
	modelConfig.MergeOptions(opts)

	// Check if we have dynamic backend options
	if opts.DynamicBackendOptions != nil {
		// Create a dynamic backend configuration
		dynamicConfig := config.BackendConfig{
			ID:          "dynamic",
			BackendType: opts.DynamicBackendOptions.Type,
			Config:      json.RawMessage(fmt.Sprintf(`{"base_url": "%s", "api_key": "%s"}`, opts.DynamicBackendOptions.URL, opts.DynamicBackendOptions.Key)),
			Options:     config.BackendOptions{},
		}

		// Create the dynamic backend
		dynamicBackend, err := createBackend(dynamicConfig)
		if err != nil {
			return modelConfig, nil, err
		}

		return modelConfig, dynamicBackend, nil
	}

	// Get the right backend
	backend, err := b.getBackend(modelConfig.Backend)
	if err != nil {
		return modelConfig, nil, err
	}

	return modelConfig, backend, nil
}

func processInferenceResponse(resp types.InferenceResponse, modelConfig config.ModelConfig) types.InferenceResponse {
	// Update the model name in the response
	resp.Model = modelConfig.GetID()

	// Generate an id if it doesn't already exist
	if resp.ID == "" {
		resp.ID = "inferix" // TODO: Generate a unique id
	}

	// Inject a default finish reason if it doesn't exist
	if resp.Done && resp.Response.FinishReason == "" {
		resp.Response.FinishReason = types.FinishReason_Stop
	}

	return resp
}

func (b *Backends) getBackend(name string) (types.Backend, error) {
	backend, ok := b.backends[name]
	if !ok {
		return nil, fmt.Errorf("backend %s not found", name)
	}
	return backend, nil
}

func injectFnCall(messages []types.InferenceMessage, functions []types.Tool) []types.InferenceMessage {
	var content strings.Builder
	content.WriteString("You may use the following FUNCTIONS in the response. Only use one function at a time. Give output in following OUTPUT_FORMAT if you want to call a function.\n\nFUNCTIONS:\n")

	for _, f := range functions {
		content.WriteString(fmt.Sprintf("- Name: %s\n", f.Name))
		if f.Description != "" {
			content.WriteString(fmt.Sprintf("  Description: %s\n", f.Description))
		}
		argsJson, _ := json.Marshal(f.Args)
		content.WriteString(fmt.Sprintf("  Parameter JSON Schema: %s\n\n", string(argsJson)))
	}

	content.WriteString(`\n\nOUTPUT_FORMAT:
Parameter Selection:
<Provide the step by step thought process to select the parameters. Go through the entire conversation>

Function Call:
<code>
{
    "type": "FUNC_CALL",
    "reasoning": "<reasoning for choosing the parameters>",
    "name": "<name of function>",
    "parameters": "<value to pass to function as parameter>"
}
</code>`)

	// Append a new system message to the messages slice
	newMessage := types.InferenceMessage{
		Role:    "system",
		Content: content.String(),
	}
	return append(messages, newMessage)
}

// createBackend creates a new backend based on the provided configuration
func createBackend(config config.BackendConfig) (types.Backend, error) {
	var backend types.Backend
	var err error

	switch config.BackendType {
	case "openai":
		backend, err = openai.New(config)
	case "ollama":
		backend, err = ollama.New(config)
	case "tei":
		backend, err = tei.New(config)
	default:
		return nil, fmt.Errorf("unsupported backend type: %s", config.BackendType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create backend: %v", err)
	}

	return backend, nil
}
