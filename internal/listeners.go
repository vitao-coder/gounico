package internal

import (
	"gounico/config"
	"gounico/pkg/listener"
	"gounico/pkg/listener/domain"
	"gounico/pkg/listener/service"
	"gounico/pkg/logging"
	"gounico/pkg/messaging"
	"gounico/pkg/worker"
	"gounico/pkg/worker/workerPool"
	"net/http"

	"go.uber.org/fx"
)

var ListenersModule = fx.Provide(
	NewPostListenerFeira,
	NewWorkerPool,
	NewListenerService)

type Listeners struct {
	ListenersContracts []service.PostListener `group:"listeners"`
	fx.In
}

type ListenerPostOutput struct {
	fx.Out
	ListenerContract []service.PostListener `group:"listeners,flatten"`
}

func NewPostListenerFeira(configuration config.Configuration) ListenerPostOutput {
	output := ListenerPostOutput{
		ListenerContract: []service.PostListener{},
	}
	for _, configFeiraListener := range configuration.Messaging.Configurations {
		listenerFeira := domain.NewGenericListener(configFeiraListener.Topic, configFeiraListener.Subscriber, configFeiraListener.URL, false)
		output.ListenerContract = append(output.ListenerContract, listenerFeira)
	}
	return output
}

func NewWorkerPool(configuration config.Configuration) worker.Worker {
	workerService := workerPool.NewWorkerPool(configuration.Messaging.WorkerPoolLimit)
	return workerService
}

func NewListenerService(configuration config.Configuration, worker worker.Worker, logger logging.Logger, httpClient *http.Client, pulsar messaging.Messaging, postListener Listeners) listener.Listener {
	listenerService := service.NewListener(worker, logger, httpClient, pulsar, configuration.Messaging.ConsumerLimit, postListener.ListenersContracts...)
	return listenerService
}
