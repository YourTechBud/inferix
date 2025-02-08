package tei

// RunFnInjection is always false for the TEI backend
func (backend *TEI) RunFnInjection() bool {
	return false
}

// EnableEmbeddingsAPI returns true if the backend supports embeddings
func (backend *TEI) EnableEmbeddingsAPI() bool {
	return true
}
