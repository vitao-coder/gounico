package novafeira

import (
	"context"
	"fmt"
	"gounico/pkg/errors"
	"gounico/pkg/logging"
	"gounico/pkg/messaging/pulsar"
	"gounico/pkg/render"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	pulsarApache "github.com/apache/pulsar-client-go/pulsar"
)

type FeiraPublisher struct {
	publisherClient pulsar.PulsarClient
	logger          logging.Logger
}

func NovaFeiraPublisher(publisherClient pulsar.PulsarClient, logger logging.Logger) FeiraPublisher {
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.RenderRequestError(w, err)
		return
	}
	exists, producer := h.publisherClient.ExistsGetProducer("feiraLivre")
	var apiError errors.ServiceError
	if !exists {
		apiError = *errors.NotFoundError()
		render.RenderApiError(w, apiError)
		return
	}
	asyncMsg := &pulsarApache.ProducerMessage{
		Key:     uuid.New().String(),
		Payload: body,
	}
	ctx := context.Background()

	producer.Producer.SendAsync(ctx, asyncMsg, func(msgId pulsarApache.MessageID, msg *pulsarApache.ProducerMessage, err error) {
		if err != nil {
			h.logger.Error(ctx, fmt.Sprintf("%s - ASYNC Error Message: %v", producer.Producer.Topic()), string(body), err)
		}
	})

	render.RenderSuccess(w, http.StatusOK, nil)
	return
}
