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
	ctx, traceSpan := openTelemetry.NewSpan(r.Context(), "Consumer - Handler")
	defer traceSpan.End()

	newFeira := &domains.FeiraRequest{}

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &newFeira)
	if err != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
		openTelemetry.AddSpanError(traceSpan, err)
		render.RenderRequestError(w, err)
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
