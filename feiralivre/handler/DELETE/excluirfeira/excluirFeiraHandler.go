package excluirfeira

import (
	"gounico/feiralivre"
	"gounico/pkg/render"
	"net/http"

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
	return "/excluirfeira/distrito/{iddistrito}/feira/{id}"
}

func (h ExcluirFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idDistrito := chi.URLParam(r, "iddistrito")
	id := chi.URLParam(r, "id")

	apiError := h.feiraLivreService.ExcluirFeira(r.Context(), id, idDistrito)
	if apiError != nil {
		render.RenderApiError(w, *apiError)
		return
	}
	render.RenderStatusCode(w, http.StatusNoContent)
}
