package ollama

// RunFnInjection returns true if the backend supports function injection
func (backend *Ollama) RunFnInjection() bool {
	return backend.options.InjectFnCallPrompt
}
