package alterarfeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domain"
	"io/ioutil"
	"net/http"
)

type AlteraFeiraHandler struct {
	feiraLivreService feiralivre.FeiraLivre
}

func NewAlteraFeiraHandler(feiraLivreService feiralivre.FeiraLivre) AlteraFeiraHandler {
	return AlteraFeiraHandler{
		feiraLivreService: feiraLivreService,
	}
}

func (h AlteraFeiraHandler) HttpMethod() string {
	return "PUT"
}

func (h AlteraFeiraHandler) HttpPath() string {
	return "/alterafeira/{id}"
}

func (h AlteraFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newFeira := &domain.FeiraRequest{}

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &newFeira)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error unmarshal request.")
		return
	}

	apiError := h.feiraLivreService.NovaFeira(r.Context(), newFeira)
	if apiError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiError.HttpStatusCode)
		json.NewEncoder(w).Encode(apiError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
