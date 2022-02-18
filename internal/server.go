package internal

import (
	"context"
	"gounico/config"
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

func NewServer(logger logging.Logger, endpointsRouter Router) *chi.Mux {
	logger.Info(context.Background(), "Starting registering endpoints in server...", nil)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(timeout * time.Second))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	logger.Info(context.Background(), "Server endpoints registered...", nil)
	return r
}

func StartServer(lc fx.Lifecycle, logger logging.Logger, server *chi.Mux, config config.Configuration) {
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

func ListenAndServe() {
	ServerModule := fx.Provide(
		NewConfig,
		NewServer,
	)
	app := fx.New(fx.Options(
		PackagesModule,
		ServerModule,
		//RepositoryModule,
		ServicesModule,
	), fx.Invoke(StartServer, StartRepository))
	app.Run()
}
