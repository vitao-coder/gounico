package buscarfeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domain"
	"gounico/pkg/errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BuscarFeiraHandler struct {
	feiraLivreService feiralivre.FeiraLivre
}

func NewBuscarFeiraHandler(feiraLivreService feiralivre.FeiraLivre) BuscarFeiraHandler {
	return BuscarFeiraHandler{
		feiraLivreService: feiraLivreService,
	}
}

func (h BuscarFeiraHandler) HttpMethod() string {
	return "GET"
}

func (h BuscarFeiraHandler) HttpPath() string {
	return "/buscarfeira/{buscapor}/{param}"
}

func (h BuscarFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buscapor := chi.URLParam(r, "buscapor")
	param := chi.URLParam(r, "param")

	var feirasRetornadas []*domain.Feira
	var apiError *errors.ServiceError
	switch buscapor {
	case "bairro":
		feirasRetornadas, apiError = h.feiraLivreService.BuscarFeiraPorBairro(r.Context(), param)
	case "distrito":
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error in URL params.")
		return
	}

	if apiError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiError.HttpStatusCode)
		json.NewEncoder(w).Encode(apiError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feirasRetornadas)
}
