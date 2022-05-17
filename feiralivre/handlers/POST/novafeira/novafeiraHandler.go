package novafeira

import (
	"encoding/json"
	"gounico/feiralivre"
	"gounico/feiralivre/domains"
	"gounico/pkg/render"
	"gounico/pkg/telemetry"
	"io/ioutil"
	"net/http"
)

type NovaFeiraHandler struct {
	feiraLivreService feiralivre.FeiraLivre
	telemetry         telemetry.OpenTelemetry
}

func NewNovaFeiraHandler(feiraLivreService feiralivre.FeiraLivre, telemetry telemetry.OpenTelemetry) NovaFeiraHandler {
	return NovaFeiraHandler{
		feiraLivreService: feiraLivreService,
		telemetry:         telemetry,
	}
}

func (h NovaFeiraHandler) HttpMethod() string {
	return "POST"
}

func (h NovaFeiraHandler) HttpPath() string {
	return "/novafeira"
}

func (h NovaFeiraHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.telemetry.Start(r.Context(), "NewNovaFeiraHandler")
	defer h.telemetry.End()
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
	h.telemetry.SuccessSpan("Success generated")
	render.RenderSuccess(w, http.StatusCreated, nil)
	return
}
