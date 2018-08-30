package client

import (
	"fmt"
)

type AmbariError struct {
	Code    int
	Message string
}

func (e AmbariError) Error() string {
	return e.Message
}

func NewAmbariError(code int, message string, params ...interface{}) AmbariError {
	return AmbariError{
		Code:    code,
		Message: fmt.Sprintf(message, params...),
	}
}
