package buscarfeira

import (
	"gounico/feiralivre"
	"gounico/internal/render"
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
		render.RenderApiError(w, *apiError)
		return
	}

	render.RenderSuccess(w, http.StatusOK, feirasRetornadas)
}
