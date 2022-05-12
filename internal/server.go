package internal

import (
	"context"
	"fmt"
	"gounico/config"
	"gounico/internal/listener"
	"gounico/pkg/logging"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
)

const (
	timeout = 60
)

type HTTPEndpoint interface {
	http.Handler
	HttpMethod() string
	HttpPath() string
}

type HTTPConsumer interface {
	http.Handler
	HttpMethod() string
	HttpPath() string
}

type MuxRouter struct {
	*chi.Mux
}

type MuxConsumer struct {
	*chi.Mux
}

type Consumer struct {
	Consumers []HTTPConsumer `group:"consumers"`
	fx.In
}

type Router struct {
	Endpoints []HTTPEndpoint `group:"endpoints"`
	fx.In
}

func NewConfig() config.Configuration {
	absPath, _ := filepath.Abs("../gounico/config/config.yaml")
	f, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg config.Configuration
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewServer(logger logging.Logger, endpointsRouter Router) *MuxRouter {
	logger.Info(context.Background(), "Starting registering endpoints in server...", nil)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout * time.Second))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	for _, endpoint := range endpointsRouter.Endpoints {
		r.Method(endpoint.HttpMethod(), endpoint.HttpPath(), endpoint)
		logger.Info(context.Background(), fmt.Sprintf("Method: %s - Pattern: %s - Registered.", endpoint.HttpMethod(), endpoint.HttpPath()), nil)
	}

	logger.Info(context.Background(), "Server endpoints registered...", nil)
	return &MuxRouter{r}
}

func NewConsumer(logger logging.Logger, endpointsConsumers Consumer) *MuxConsumer {
	logger.Info(context.Background(), "Starting registering consumers endpoints in server...", nil)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout * time.Second))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	for _, endpoint := range endpointsConsumers.Consumers {
		r.Method(endpoint.HttpMethod(), endpoint.HttpPath(), endpoint)
		logger.Info(context.Background(), fmt.Sprintf("Consumer: %s - Pattern: %s - Registered.", endpoint.HttpMethod(), endpoint.HttpPath()), nil)
	}

	logger.Info(context.Background(), "Consumer endpoints registered...", nil)
	return &MuxConsumer{r}
}

func StartServer(lc fx.Lifecycle, logger logging.Logger, server *MuxRouter, config config.Configuration) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(ctx, "Start server", nil)
			go http.ListenAndServe(":"+config.Server.Port, server)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(ctx, "Stop server", nil)
			return nil
		},
	})
}

func StartListener(lc fx.Lifecycle, logger logging.Logger, listenerConsumers *MuxConsumer, config config.Configuration, listenerService listener.Listener) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(ctx, "Start listener / consumers", nil)
			go http.ListenAndServe(":"+config.Worker.Port, listenerConsumers)
			go listenerService.RunListenerService()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(ctx, "Stop listener / consumers", nil)
			listenerService.StopService()
			ctx.Done()
			return nil
		},
	})
}

func ListenAndServe() {
	ServerModule := fx.Provide(
		NewConfig,
		NewServer,
		NewConsumer,
	)
	app := fx.New(fx.Options(
		PackagesModule,
		ServerModule,
		RepositoryModule,
		ServicesModule,
		ConsumersModule,
		ListenersModule,
		HandlersModule,
	), fx.Invoke(StartServer, StartListener))
	app.Run()
}
