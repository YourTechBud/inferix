package openai

// RunFnInjection returns true if the backend supports function injection
func (backend *OpenAI) RunFnInjection() bool {
	return backend.options.InjectFnCallPrompt
}

// EnableEmbeddingsAPI returns true if the backend supports embeddings
func (backend *OpenAI) EnableEmbeddingsAPI() bool {
	return backend.options.EnableEmbeddingsAPI
}
