package processcsv

import (
	"encoding/json"
	"gounico/loaddata"
	"io"
	"net/http"
)

type ProcessCSVHandler struct {
	loadDataService loaddata.LoadData
}

func NewProcessCSVHandler(loadDataService loaddata.LoadData) ProcessCSVHandler {
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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Bad Request")
		return
	}

	err = h.loadDataService.ProcessCSVToDatabase(bodyBytes)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error on process CSV.")
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
