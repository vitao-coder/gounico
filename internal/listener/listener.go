package listener

import (
	"bytes"
	"context"
	"fmt"
	"gounico/constants"
	"gounico/pkg/apiclient"
	"gounico/pkg/logging"
	"gounico/pkg/messaging/pulsar/client"
	"gounico/pkg/worker"
	"net/http"
	"sync"

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
	pulsarClient         client.PulsarClient
	consumerChannelLimit int
}

func NewListener(workerService worker.Worker,
	logger logging.Logger,
	httpClient *http.Client,
	pulsarClient client.PulsarClient, consumerChannelLimit int, postListeners ...PostListener) *Listener {
	return &Listener{
		workerService:        workerService,
		logger:               logger,
		httpClient:           httpClient,
		postListeners:        postListeners,
		pulsarClient:         pulsarClient,
		consumerChannelLimit: consumerChannelLimit,
	}
}

func (l *Listener) RunListenerService(ctx context.Context) {
	l.ctx = ctx
	go l.workerService.Run(ctx)
	l.logger.Info(ctx, "Started listener - Worker service", nil)
	go l.listenResults()
	l.logger.Info(ctx, "Started listener - Listen service", nil)
	go l.runConsumerChannels(ctx)
	l.logger.Info(ctx, "Started listener - Consumer service", nil)
}

func (l *Listener) listenResults() {
	for {
		select {
		case r, ok := <-l.workerService.Results():
			if !ok {
				continue
			}
			if r.Error != nil {
				l.logger.Error(l.ctx, fmt.Sprintf("ERROR - Worker: %s", r.WorkerJobDescriptor), r.Result, r.Error)
			} else {
				l.logger.Info(l.ctx, fmt.Sprintf("SUCCESS - Worker: %s", r.WorkerJobDescriptor), r.Result)
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
		go l.startConsumingMessages(ctx, listener.Subscription(), l.logger, &wg, channel, consumer, listener.URL())
	}
}

func (l *Listener) startConsumingMessages(ctx context.Context, consumerName string, logger logging.Logger, wg *sync.WaitGroup, messages <-chan pulsar.ConsumerMessage, consumer pulsar.Consumer, postURL string) {
	defer wg.Done()

	for {
		select {
		case chMsg, ok := <-messages:
			if !ok {
				return
			}
			msg := chMsg.Message
			logger.Info(ctx, fmt.Sprintf("%s - Received message in consumer. Message: %v", consumerName, msg.Payload()), string(msg.Payload()))
			l.addWorkerConsumerPOST(consumerName, postURL, msg)
			consumer.Ack(chMsg)
		case <-ctx.Done():
			logger.Error(ctx, fmt.Sprintf("%s - Cancelled consumer. Error detail: %v", consumerName, ctx.Err()), nil, ctx.Err())
			return
		}
	}
}

func (l *Listener) postFunction() worker.JobFunction {
	var postFunction = func(ctx context.Context, params ...interface{}) (interface{}, error) {
		consumerName := params[0].(string)
		msg := params[1].(pulsar.Message)
		urlToPost := params[2].(string)
		payloadToPost := msg.Payload()

		r := bytes.NewReader(payloadToPost)
		result, err := apiclient.Post(ctx, l.httpClient, urlToPost, constants.ContentTypeJson, r)
		if err != nil {
			l.logger.Error(ctx, fmt.Sprintf("%s - Consumer message send ERROR.", consumerName), msg.Properties(), err)
			return nil, err
		}
		l.logger.Info(ctx, fmt.Sprintf("%s - Consumer message send SUCCESS.", consumerName), msg.Properties())
		return result, nil
	}
	return postFunction
}

func (l *Listener) addWorkerConsumerPOST(name string, urlToPost string, msg pulsar.Message) worker.WorkerJob {
	var params []interface{}
	params = append(params, name)
	params = append(params, urlToPost)
	params = append(params, msg)
	workerJob := worker.NewWorkerJob(name, l.postFunction(), params)
	return workerJob
}
