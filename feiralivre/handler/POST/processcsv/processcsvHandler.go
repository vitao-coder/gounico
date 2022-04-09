package processcsv

import (
	"bufio"
	"encoding/json"
	"gounico/feiralivre"
	"io"
	"net/http"
)

type ProcessCSVHandler struct {
	loadDataService feiralivre.ProcessCSV
}

func NewProcessCSVHandler(loadDataService feiralivre.ProcessCSV) ProcessCSVHandler {
	return ProcessCSVHandler{
		loadDataService: loadDataService,
	}
}

func (h ProcessCSVHandler) HttpMethod() string {
	return "POST"
}

func (h ProcessCSVHandler) HttpPath() string {
	return "/csvprocessor"
}

func (h ProcessCSVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error receiving file bytes.")
		return
	}
	bodyBytes := bufio.NewReader(file)
	bytesBody, err := io.ReadAll(bodyBytes)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error receiving file bytes.")
		return
	}

	apiErr := h.loadDataService.ProcessCSVToDatabase(r.Context(), bytesBody)
	if apiErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.HttpStatusCode)
		json.NewEncoder(w).Encode(apiErr)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
