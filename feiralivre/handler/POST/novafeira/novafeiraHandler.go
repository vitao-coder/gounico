package novafeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domain"
	"io/ioutil"
	"net/http"
)

type NovaFeiraHandler struct {
	feiraLivreService feiralivre.FeiraLivre
}

func NewNovaFeiraHandler(feiraLivreService feiralivre.FeiraLivre) NovaFeiraHandler {
	return NovaFeiraHandler{
		feiraLivreService: feiraLivreService,
	}
}

func (h NovaFeiraHandler) HttpMethod() string {
	return "POST"
}

func (h NovaFeiraHandler) HttpPath() string {
	return "/novafeira"
}

func (h NovaFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newFeira := &domain.FeiraRequest{}

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &newFeira)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error unmarshal request.")
		return
	}

	apiError := h.feiraLivreService.SaveFeira(r.Context(), newFeira)
	if apiError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiError.HttpStatusCode)
		json.NewEncoder(w).Encode(apiError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
