package novafeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domain"
	"gounico/internal/render"
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
		render.RenderRequestError(w, err)
		return
	}

	apiError := h.feiraLivreService.SaveFeira(r.Context(), newFeira)
	if apiError != nil {
		render.RenderApiError(w, *apiError)
		return
	}

	render.RenderSuccess(w, http.StatusCreated, nil)
	return
}
