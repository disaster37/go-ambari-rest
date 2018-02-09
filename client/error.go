package client

type AmbariError struct {
	Code    int
	Message string
}

func (e AmbariError) Error() string {
	return e.Message
}

func NewAmbariError(code int, message string) AmbariError {
	return AmbariError{
		Code:    code,
		Message: message,
	}
}
