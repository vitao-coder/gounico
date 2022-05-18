package novafeira

import (
	"encoding/json"
	"fmt"
	"gounico/feiralivre"
	"gounico/feiralivre/domains"
	"gounico/pkg/render"
	"gounico/pkg/telemetry/openTelemetry"
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
	ctx, traceSpan := openTelemetry.NewSpan(r.Context(), h.HttpPath()+" - Handler")

	defer traceSpan.End()

	newFeira := &domains.FeiraRequest{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
		openTelemetry.AddSpanError(traceSpan, err)
		render.RenderRequestError(w, err)
		return
	}
	errMarshal := json.Unmarshal(body, &newFeira)
	if errMarshal != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", errMarshal.Error()))
		openTelemetry.AddSpanError(traceSpan, errMarshal)
		render.RenderRequestError(w, errMarshal)
		return
	}

	var mapsTags map[string]string
	json.Unmarshal(body, &mapsTags)
	openTelemetry.AddSpanTags(traceSpan, mapsTags)
	apiError := h.feiraLivreService.SaveFeira(ctx, newFeira)
	if apiError != nil {
		openTelemetry.FailSpan(traceSpan, apiError.Error())
		openTelemetry.AddSpanError(traceSpan, apiError)
		render.RenderApiError(w, *apiError)
		return
	}
	openTelemetry.SuccessSpan(traceSpan, fmt.Sprintf("StatusCode: %d", http.StatusCreated))
	render.RenderSuccess(w, http.StatusCreated, nil)
	return
}
