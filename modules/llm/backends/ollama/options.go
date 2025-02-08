package ollama

// RunFnInjection returns true if the backend supports function injection
func (backend *Ollama) RunFnInjection() bool {
	return backend.options.InjectFnCallPrompt
}

// EnableEmbeddingsAPI returns true if the backend supports embeddings
func (backend *Ollama) EnableEmbeddingsAPI() bool {
	return backend.options.EnableEmbeddingsAPI
}
