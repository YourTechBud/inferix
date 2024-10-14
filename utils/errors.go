package utils

type StandardError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

func NewStandardError(status int, message, errorType string) StandardError {
	return StandardError{
		Status:  status,
		Message: message,
		Type:    errorType,
	}
}

func (e StandardError) Error() string {
	return e.Message
}

var _ error = (*StandardError)(nil)
