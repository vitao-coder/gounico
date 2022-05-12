package internal

import (
	"gounico/config"
	"gounico/internal/listener"
	"gounico/internal/listener/domain"
	"gounico/internal/listener/service"
	"gounico/pkg/logging"
	"gounico/pkg/messaging/pulsar"
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
	ListenerContract service.PostListener `group:"listeners"`
}

func NewPostListenerFeira(configuration config.Configuration) ListenerPostOutput {
	configFeiraListener := configuration.Messaging.Configurations[0]
	listenerFeira := domain.NewGenericListener(configFeiraListener.Topic, configFeiraListener.Subscriber, configFeiraListener.URL, false)
	return ListenerPostOutput{
		ListenerContract: listenerFeira,
	}

}

func NewWorkerPool(configuration config.Configuration) worker.Worker {
	workerService := workerPool.NewWorkerPool(configuration.Messaging.ConsumerLimit)
	return workerService
}

func NewListenerService(configuration config.Configuration, worker worker.Worker, logger logging.Logger, httpClient *http.Client, pulsar pulsar.PulsarClient, postListener Listeners) listener.Listener {
	listenerService := service.NewListener(worker, logger, httpClient, pulsar, configuration.Messaging.ConsumerLimit, postListener.ListenersContracts...)
	return listenerService
}
