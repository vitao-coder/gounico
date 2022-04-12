package excluirfeira

import (
	"gounico/feiralivre"
	"gounico/internal/render"
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
	return "/excluirfeira/{id}"
}

func (h ExcluirFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")

	apiError := h.feiraLivreService.ExcluirFeira(r.Context(), param)

	if apiError != nil {
		render.RenderApiError(w, *apiError)
		return
	}
	render.RenderSuccess(w, http.StatusNoContent, nil)
}
