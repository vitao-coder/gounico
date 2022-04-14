package alterarfeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domain"
	"gounico/pkg/render"
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
	updFeira := &domain.FeiraRequest{}

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &updFeira)
	if err != nil {
		render.RenderRequestError(w, err)
		return
	}

	apiError := h.feiraLivreService.SaveFeira(r.Context(), updFeira)
	if apiError != nil {
		render.RenderApiError(w, *apiError)
		return
	}
	render.RenderSuccess(w, http.StatusOK, nil)
	return
}
