package api

import (
	"context"
	"net/http"
	"time"

	"github.com/goawwer/admintasks/config"
	"github.com/goawwer/admintasks/store"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type API struct {
	config  *config.Config
	router  *mux.Router
	logger  *logrus.Logger
	storage *store.Storage
}

func New(cfg *config.Config) *API {
	return &API{
		config: cfg,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
}

func (api *API) Start(ctx context.Context) error {
	if err := api.ConfigureLogger(); err != nil {
		return err
	}

	if err := api.ConfigureStorage(); err != nil {
		return err
	}

	api.ConfigureRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Или точные origin'ы, напр: "http://localhost:3000"
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:    api.config.BindAddr,
		Handler: c.Handler(api.router),
	}

	api.logger.WithFields(logrus.Fields{
		"addr": api.config.BindAddr,
	}).Info("Starting API server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			api.logger.WithError(err).Error("HTTP server crashed")
		}
	}()

	<-ctx.Done()

	api.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		api.logger.WithError(err).Error("Graceful shutdown failed")
	}

	return nil
}
