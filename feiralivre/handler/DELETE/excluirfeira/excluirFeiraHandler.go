package excluirfeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ExcluirFeiraHandler struct {
	feiraLivreService feiralivre.FeiraLivre
}

func NewExcluirFeiraHandler(feiraLivreService feiralivre.FeiraLivre) ExcluirFeiraHandler {
	return ExcluirFeiraHandler{
		feiraLivreService: feiraLivreService,
	}
}

func (h ExcluirFeiraHandler) HttpMethod() string {
	return "DELETE"
}

func (h ExcluirFeiraHandler) HttpPath() string {
	return "/excluirfeira/{id}"
}

func (h ExcluirFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error in URL params.")
		return
	}

	apiError := h.feiraLivreService.ExcluirFeira(r.Context(), uint(id))

	if apiError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiError.HttpStatusCode)
		json.NewEncoder(w).Encode(apiError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(nil)
}
