package novafeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domains"
	"gounico/pkg/render"
	"io/ioutil"
	"net/http"
)

type NovaFeiraConsumer struct {
	feiraLivreService feiralivre.FeiraLivre
}

func NewNovaFeiraConsumer(feiraLivreService feiralivre.FeiraLivre) NovaFeiraConsumer {
	return NovaFeiraConsumer{
		feiraLivreService: feiraLivreService,
	}
}

func (h NovaFeiraConsumer) HttpMethod() string {
	return "POST"
}

func (h NovaFeiraConsumer) HttpPath() string {
	return "/consumers/novafeira"
}

func (h NovaFeiraConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newFeira := &domains.FeiraRequest{}

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
