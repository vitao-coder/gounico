package novafeira

import (
	"encoding/json"
	"fmt"
	"gounico/feiralivre/domains"
	"gounico/pkg/errors"
	"gounico/pkg/logging"
	"gounico/pkg/messaging"
	"gounico/pkg/render"
	"gounico/pkg/telemetry/openTelemetry"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	pulsarApache "github.com/apache/pulsar-client-go/pulsar"
)

type FeiraPublisher struct {
	publisherClient messaging.Messaging
	logger          logging.Logger
}

func NovaFeiraPublisher(publisherClient messaging.Messaging, logger logging.Logger) FeiraPublisher {
	return FeiraPublisher{
		publisherClient: publisherClient,
		logger:          logger,
	}
}

func (h FeiraPublisher) HttpMethod() string {
	return "POST"
}

func (h FeiraPublisher) HttpPath() string {
	return "/publisher/novafeira"
}

func (h FeiraPublisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceSpan := openTelemetry.SpanFromContext(ctx)

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
		openTelemetry.AddSpanError(traceSpan, errMarshal)
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", errMarshal.Error()))
		render.RenderRequestError(w, errMarshal)
		return
	}
	exists, producer := h.publisherClient.ExistsGetProducer("feiraLivre")
	var apiError errors.ServiceError
	if !exists {
		apiError = *errors.NotFoundError()
		openTelemetry.AddSpanError(traceSpan, apiError)
		render.RenderApiError(w, apiError)
		return
	}

	var mapsTags map[string]string
	json.Unmarshal(body, &mapsTags)
	openTelemetry.AddSpanTags(traceSpan, mapsTags)

	asyncMsg := &pulsarApache.ProducerMessage{
		Key:     uuid.New().String(),
		Payload: body,
	}

	ctxPublisher := openTelemetry.BuildAndInjectSpanOnMessageContext(ctx, "Server - pulsar.Producer", asyncMsg)

	producer.Producer.SendAsync(ctxPublisher, asyncMsg, func(msgId pulsarApache.MessageID, msg *pulsarApache.ProducerMessage, err error) {
		traceProducerSpan := openTelemetry.SpanFromContext(ctxPublisher)
		defer traceProducerSpan.End()
		if err != nil {
			h.logger.Error(ctx, fmt.Sprintf("%s - ASYNC Error Message: %v", producer.Producer.Topic()), string(body), err)
			openTelemetry.FailSpan(traceSpan, err.Error())
			openTelemetry.AddSpanError(traceSpan, err)
		}
	})
	render.RenderSuccess(w, http.StatusOK, nil)
	return
}
