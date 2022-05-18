package errors

import (
	"fmt"
	"net/http"
)

type ServiceError struct {
	HttpStatusCode int      `json:"http_status_code,omitempty"`
	Message        string   `json:"message,omitempty"`
	Causes         []string `json:"causes,omitempty"`
}

func (s ServiceError) Error() string {
	return fmt.Sprintf("StatusCode: %d - Message: %s", s.HttpStatusCode, s.Message)
}

func InternalServerError(message string, err error) *ServiceError {
	return &ServiceError{
		HttpStatusCode: http.StatusInternalServerError,
		Message:        message,
		Causes:         []string{"internal_server_error", err.Error()},
	}
}

func BadRequestError(message string) *ServiceError {
	return GenericError(message, http.StatusBadRequest, []string{"bad_request"})
}

func NotFoundError() *ServiceError {
	return GenericError("Not Found", http.StatusNotFound, []string{"not_found"})

}

func GenericError(message string, statusCode int, causes []string) *ServiceError {
	return &ServiceError{
		HttpStatusCode: statusCode,
		Message:        message,
		Causes:         causes,
	}
}
