package render

import (
	"encoding/json"
	"gounico/pkg/errors"
	"net/http"
)

func renderApiResponse(w http.ResponseWriter, renderObject interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	err := json.NewEncoder(w).Encode(renderObject)
	if err != nil {
		panic(err)
	}
}

func RenderRequestError(w http.ResponseWriter, error error) {
	var serviceError = errors.ServiceError{
		HttpStatusCode: http.StatusBadRequest,
		Message:        error.Error(),
		Causes:         []string{"bad_request"},
	}
	RenderApiError(w, serviceError)
}

func RenderApiError(w http.ResponseWriter, error errors.ServiceError) {
	renderApiResponse(w, error, error.HttpStatusCode)
}

func RenderSuccess(w http.ResponseWriter, httpStatusCode int, objectReturn interface{}) {
	renderApiResponse(w, objectReturn, httpStatusCode)
}
