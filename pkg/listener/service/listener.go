package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gounico/global"
	"gounico/pkg/apiclient"
	"gounico/pkg/logging"
	pulsarMsg "gounico/pkg/messaging"
	"gounico/pkg/telemetry/openTelemetry"
	"gounico/pkg/worker"
	"gounico/pkg/worker/domain"
	"io/ioutil"
	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
)

type PostListener interface {
	Topic() string
	Subscription() string
	URL() string
	IsShared() bool
}

var subscriptionsConsumerLimit = 2

type Listener struct {
	workerService        worker.Worker
	logger               logging.Logger
	ctx                  context.Context
	httpClient           *http.Client
	postListeners        []PostListener
	pulsarClient         pulsarMsg.Messaging
	consumerChannelLimit int
}

func NewListener(workerService worker.Worker,
	logger logging.Logger,
	httpClient *http.Client,
	pulsarClient pulsarMsg.Messaging, consumerChannelLimit int, postListeners ...PostListener) *Listener {
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
	l.createConsumerChannels(l.ctx)
	l.logger.Info(l.ctx, "Started listener - Consumers services", nil)
}

func (l *Listener) StopService() {
	l.ctx.Done()
}

func (l *Listener) listenResults() {
	for {
		traceSpan := openTelemetry.SpanFromContext(l.ctx)
		select {
		case r, ok := <-l.workerService.Results():
			if !ok {
				continue
			}
			if r.Error != nil {
				openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", r.Error.Error()))
				openTelemetry.AddSpanError(traceSpan, r.Error)
				l.logger.Error(l.ctx, fmt.Sprintf("%s - ERROR - Worker Job", r.WorkerJobDescriptor), r.Result, r.Error)
			}
			openTelemetry.SuccessSpan(traceSpan, "Success")
		case <-l.workerService.Finished():
			return
		default:
		}
	}
}

func (l *Listener) createConsumerChannels(ctx context.Context) error {
	for _, postListener := range l.postListeners {
		for i := 0; i < subscriptionsConsumerLimit; i++ {
			name := fmt.Sprintf("%s-%d", postListener.Subscription(), i)
			err := l.pulsarClient.CreateConsumerWithChannels(postListener.Topic(), postListener.Subscription(), name, l.consumerChannelLimit)
			if err != nil {
				return err
			}
			l.runConsumerChannel(ctx, name, postListener)
		}
	}
	return nil
}

func (l *Listener) runConsumerChannel(ctx context.Context, consumerName string, listener PostListener) {
	consumer, channel := l.pulsarClient.GetConsumer(listener.Topic(), listener.Subscription(), consumerName)
	go startConsumingMessages(ctx, consumerName, l.logger, channel, consumer, listener.URL(), l.workerService, l.httpClient)
}

func startConsumingMessages(ctx context.Context, consumerName string, logger logging.Logger, messages <-chan pulsar.ConsumerMessage, consumer pulsar.Consumer, postURL string, workerService worker.Worker, httpClient *http.Client) {
	for {
		select {
		case chMsg, ok := <-messages:
			if !ok {
				continue
			}
			ctxMsg := openTelemetry.ExtractTraceContextFromMessage(ctx, chMsg)
			ctxConsumer, traceSpan := openTelemetry.TraceContextSpan(ctxMsg, "Listener - pulsar.Consumer")
			msg := chMsg.Message
			json, _ := json.Marshal(msg.Payload())
			logger.Debug(ctxConsumer, fmt.Sprintf("%s - Consuming message. KeyMsg: %s", consumerName, msg.Key()), string(json))
			workerService.AddJobs(createWorkerJOB(fmt.Sprintf("%s - ConsumerJOB. KeyMsg %s", consumerName, msg.Key()), postURL, msg, consumer, httpClient, ctxConsumer))
			openTelemetry.SuccessSpan(traceSpan, "Success")
			traceSpan.End()
		case <-ctx.Done():
			qtdMsgs := len(messages)
			if qtdMsgs > 0 {
				logger.Error(ctx, fmt.Sprintf("%s - Cancelled. Error detail: %v", consumerName, ctx.Err()), nil, ctx.Err())
			}
			continue
		}
	}
}

func postFunction(consumer pulsar.Consumer, httpClient *http.Client) domain.JobFunction {
	var postFunction = func(ctx context.Context, params []interface{}) (interface{}, error) {
		ctx, traceSpan := openTelemetry.NewSpan(ctx, "Listener - ExecuteJobFunction")
		defer traceSpan.End()

		msg, urlToPost, payloadToPost := extractParamsToPost(params)
		r := bytes.NewReader(payloadToPost)

		result, err := apiclient.Post(ctx, httpClient, urlToPost, global.ContentTypeJson, r)
		if err != nil || result.StatusCode > 400 {
			if err == nil {
				traceSpan := openTelemetry.SpanFromContext(ctx)
				openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
				openTelemetry.AddSpanError(traceSpan, err)
				consumer.Nack(msg)
				return nil, errors.New("error +400 calling post consumer.")
			}
			consumer.Nack(msg)
			return nil, err
		}

		openTelemetry.SuccessSpan(traceSpan, fmt.Sprintf("Success"))
		consumer.Ack(msg)
		defer result.Body.Close()

		bodyResult, err := ioutil.ReadAll(result.Body)
		return bodyResult, nil
	}
	return postFunction
}

func extractParamsToPost(params []interface{}) (pulsar.Message, string, []byte) {
	var insideParams []interface{}
	insideParams = params[0].([]interface{})
	msg := insideParams[0].(pulsar.Message)
	urlToPost := insideParams[1].(string)
	payloadToPost := msg.Payload()
	return msg, urlToPost, payloadToPost
}

func createWorkerJOB(name string, urlToPost string, msg pulsar.Message, consumer pulsar.Consumer, httpClient *http.Client, ctx context.Context) domain.WorkerJob {
	var params []interface{}
	params = append(params, msg)
	params = append(params, urlToPost)
	workerJob := domain.NewWorkerJob(name, postFunction(consumer, httpClient), ctx, params)
	return workerJob
}
