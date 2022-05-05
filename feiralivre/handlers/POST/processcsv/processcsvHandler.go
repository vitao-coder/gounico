package processcsv

import (
	"bufio"
	"gounico/feiralivre"
	"gounico/infrastructure/render"
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
		render.RenderRequestError(w, err)
		return
	}
	bodyBytes := bufio.NewReader(file)
	bytesBody, err := io.ReadAll(bodyBytes)
	if err != nil {
		render.RenderRequestError(w, err)
		return
	}

	apiErr := h.loadDataService.ProcessCSVToDatabase(r.Context(), bytesBody)
	if apiErr != nil {
		render.RenderApiError(w, *apiErr)
		return
	}
	render.RenderSuccess(w, http.StatusOK, nil)
	return
}
