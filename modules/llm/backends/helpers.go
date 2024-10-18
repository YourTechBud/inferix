package backends

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

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
