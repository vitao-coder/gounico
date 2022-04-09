package buscarfeira

import (
	"encoding/json"
	"gounico/feiralivre"
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
	return "/buscarfeira/{entity_type}/{entityID}"
}

func (h BuscarFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	entityID := chi.URLParam(r, "entityID")

	feirasRetornadas, apiError := h.feiraLivreService.BuscarFeiraPorDistrito(r.Context(), entityID)

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
