package tei

// RunFnInjection is always false for the TEI backend
func (backend *TEI) RunFnInjection() bool {
	return false
}
