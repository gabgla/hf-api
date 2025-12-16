package api

type APIError struct {
	Message    string `json:"message"`
	InnerError error  `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

func wrapError(newError string, original error) *APIError {
	return &APIError{
		Message:    newError,
		InnerError: original,
	}
}
