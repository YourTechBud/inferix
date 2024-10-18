package openai

// RunFnInjection returns true if the backend supports function injection
func (backend *OpenAI) RunFnInjection() bool {
	return backend.options.InjectFnCallPrompt
}
