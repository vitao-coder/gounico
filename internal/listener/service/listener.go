package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gounico/feiralivre/domains"
	"gounico/global"
	"gounico/pkg/apiclient"
	"gounico/pkg/logging"
	pulsar2 "gounico/pkg/messaging/pulsar"
	"gounico/pkg/messaging/pulsar/tracing"
	"gounico/pkg/telemetry/openTelemetry"
	"gounico/pkg/worker"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/propagation"

	"github.com/apache/pulsar-client-go/pulsar"
)

type PostListener interface {
	Topic() string
	Subscription() string
	URL() string
	IsShared() bool
}

type Listener struct {
	workerService        worker.Worker
	logger               logging.Logger
	ctx                  context.Context
	httpClient           *http.Client
	postListeners        []PostListener
	pulsarClient         pulsar2.PulsarClient
	consumerChannelLimit int
}

func NewListener(workerService worker.Worker,
	logger logging.Logger,
	httpClient *http.Client,
	pulsarClient pulsar2.PulsarClient, consumerChannelLimit int, postListeners ...PostListener) *Listener {
	return &Listener{
		workerService:        workerService,
		logger:               logger,
		httpClient:           httpClient,
		postListeners:        postListeners,
		pulsarClient:         pulsarClient,
		consumerChannelLimit: consumerChannelLimit,
	}
}

func (l *Listener) RunListenerService() {
	l.ctx = context.Background()
	go l.workerService.Run(l.ctx)
	l.logger.Info(l.ctx, "Started listener - Workers services", nil)
	go l.listenResults()
	l.logger.Info(l.ctx, "Started listener - Listeners service", nil)
	l.createConsumerChannels()
	go l.runConsumerChannels(l.ctx)
	l.logger.Info(l.ctx, "Started listener - Consumers services", nil)
	//l.TestProducerMessage()

}

func (l *Listener) StopService() {
	l.ctx.Done()
}

func (l *Listener) listenResults() {
	for {
		select {
		case r, ok := <-l.workerService.Results():
			if !ok {
				continue
			}
			if r.Error != nil {
				l.logger.Error(l.ctx, fmt.Sprintf("%s - ERROR - Worker Job", r.WorkerJobDescriptor), r.Result, r.Error)
			}
		case <-l.workerService.Finished():
			return
		default:
		}
	}
}

func (l *Listener) createConsumerChannels() error {
	for _, postListener := range l.postListeners {
		err := l.pulsarClient.CreateConsumerWithChannels(postListener.Topic(), postListener.Subscription(), l.consumerChannelLimit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Listener) runConsumerChannels(ctx context.Context) {
	var wg sync.WaitGroup
	for _, listener := range l.postListeners {
		wg.Add(1)
		consumer, channel := l.pulsarClient.GetConsumer(listener.Topic(), listener.Subscription())
		go startConsumingMessages(ctx, listener.Subscription(), l.logger, &wg, channel, consumer, listener.URL(), l.workerService, l.httpClient)
	}
}

func startConsumingMessages(ctx context.Context, consumerName string, logger logging.Logger, wg *sync.WaitGroup, messages <-chan pulsar.ConsumerMessage, consumer pulsar.Consumer, postURL string, workerService worker.Worker, httpClient *http.Client) {
	defer wg.Done()

	for {
		select {
		case chMsg, ok := <-messages:
			if !ok {
				continue
			}

			msg := chMsg.Message
			json, _ := json.Marshal(msg.Payload())
			logger.Debug(ctx, fmt.Sprintf("%s - Received message in consumer. Message: %s", consumerName, msg.Payload()), string(json))
			ctxExtracted := buildAndExtractContext(ctx, chMsg)
			workerService.AddJobs(createWorkerJOB(fmt.Sprintf("Consumer - %s MessageKey - %s", consumerName, msg.Key()), postURL, msg, consumer, httpClient, ctxExtracted))
		case <-ctx.Done():
			qtdMsgs := len(messages)
			if qtdMsgs > 0 {
				logger.Error(ctx, fmt.Sprintf("%s - Cancelled consumer. Error detail: %v", consumerName, ctx.Err()), nil, ctx.Err())
			}
			continue
		}
	}
}

func buildAndExtractContext(ctx context.Context, message pulsar.ConsumerMessage) context.Context {
	consumerAdapter := tracing.ConsumerMessageAdapter{Message: message}
	traceContext := propagation.TraceContext{}
	return traceContext.Extract(ctx, &consumerAdapter)
}

func postFunction(consumer pulsar.Consumer, httpClient *http.Client) worker.JobFunction {
	var postFunction = func(ctx context.Context, params []interface{}) (interface{}, error) {

		var insideParams []interface{}
		insideParams = params[0].([]interface{})

		msg := insideParams[0].(pulsar.Message)

		urlToPost := insideParams[1].(string)

		payloadToPost := msg.Payload()

		ctxRequest, traceSpan := openTelemetry.NewSpan(ctx, fmt.Sprintf("Listener.PostToConsumer - %s", msg.Key()))
		defer traceSpan.End()

		r := bytes.NewReader(payloadToPost)

		result, err := apiclient.Post(ctxRequest, httpClient, urlToPost, global.ContentTypeJson, r)
		if err != nil || result.StatusCode > 400 {
			if err == nil {
				return nil, errors.New("Error +400 calling POST Consumer.")
			}
			return nil, err
		}
		consumer.Ack(msg)
		defer result.Body.Close()

		bodyResult, err := ioutil.ReadAll(result.Body)
		return bodyResult, nil
	}
	return postFunction
}

func createWorkerJOB(name string, urlToPost string, msg pulsar.Message, consumer pulsar.Consumer, httpClient *http.Client, ctx context.Context) worker.WorkerJob {
	var params []interface{}
	params = append(params, msg)
	params = append(params, urlToPost)
	workerJob := worker.NewWorkerJob(name, postFunction(consumer, httpClient), ctx, params)
	return workerJob
}

func (l *Listener) TestProducerMessage() {

	exists, publisher := l.pulsarClient.ExistsGetProducer("feiraLivre")

	if !exists {
		return
	}

	ctx := context.Background()

	for i := 0; i < 10; i++ {

		requestTest := domains.FeiraRequest{
			Longitude:  "-46550164",
			Latitude:   "-23558733",
			SetCens:    "355030885000091",
			AreaP:      "3550308005040",
			CodDist:    "999",
			Distrito:   "VILA FORMOSA",
			CodSubPref: "26",
			SubPrefe:   "A-FORMOSA-CARRAO",
			Regiao5:    "Leste Novo",
			Regiao8:    "Leste 3",
			NomeFeira:  "VILA FORMOSA 2",
			Registro:   "4041-0",
			Logradouro: "RUA A",
			Numero:     "S/N",
			Bairro:     "CHUPINS",
			Referencia: "TV RUA PRETORIA A",
		}
		requestTest.Id = strconv.Itoa(i + 1)
		requestTest.CodDist = strconv.Itoa(i + 2)
		jsonRequestBytes, _ := json.Marshal(requestTest)

		asyncMsg := &pulsar.ProducerMessage{
			Payload: jsonRequestBytes,
		}

		publisher.Producer.SendAsync(ctx, asyncMsg, func(msgId pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
			l.logger.Debug(ctx, fmt.Sprintf("Send Message + IDFeira: %s / ID Distrito: %s", requestTest.Id, requestTest.Distrito), nil)
			if err != nil {
				l.logger.Error(ctx, fmt.Sprintf("%s - ASYNC Error Message: %v", publisher.Producer.Topic()), string(jsonRequestBytes), err)
			}
		})
	}
	time.Sleep(30 * time.Second)
	for i := 10; i < 20; i++ {
		requestTest := domains.FeiraRequest{
			Longitude:  "-46550164",
			Latitude:   "-23558733",
			SetCens:    "355030885000091",
			AreaP:      "3550308005040",
			CodDist:    "999",
			Distrito:   "VILA FORMOSA",
			CodSubPref: "26",
			SubPrefe:   "A-FORMOSA-CARRAO",
			Regiao5:    "Leste Novo",
			Regiao8:    "Leste 3",
			NomeFeira:  "VILA FORMOSA 2",
			Registro:   "4041-0",
			Logradouro: "RUA A",
			Numero:     "S/N",
			Bairro:     "CHUPINS",
			Referencia: "TV RUA PRETORIA A",
		}
		requestTest.Id = strconv.Itoa(i + 1)
		jsonRequestBytes, _ := json.Marshal(requestTest)

		asyncMsg := &pulsar.ProducerMessage{
			Payload: jsonRequestBytes,
		}

		publisher.Producer.SendAsync(ctx, asyncMsg, func(msgId pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
			l.logger.Debug(ctx, fmt.Sprintf("Send Message + IDFeira: %s / ID Distrito: %s", requestTest.Id, requestTest.Distrito), nil)
			if err != nil {
				l.logger.Error(ctx, fmt.Sprintf("%s - ASYNC Error Message: %v", publisher.Producer.Topic()), string(jsonRequestBytes), err)
			}
		})
	}
}
