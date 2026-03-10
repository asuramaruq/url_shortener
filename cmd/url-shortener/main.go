package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/asuramaruq/url_shortener/internal/config"
	"github.com/asuramaruq/url_shortener/internal/http-server/handlers/url/delete"
	"github.com/asuramaruq/url_shortener/internal/http-server/handlers/url/redirect"
	"github.com/asuramaruq/url_shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/asuramaruq/url_shortener/internal/http-server/middleware/logger"
	"github.com/asuramaruq/url_shortener/internal/lib/logger/sl"
	"github.com/asuramaruq/url_shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		return
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger) // default logger
	router.Use(mwLogger.New(log)) // custom middleware logger
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/save", save.New(log, storage))
	router.Get("/get/{alias}", redirect.New(log, storage))
	router.Delete("/delete/{alias}", delete.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("server started")

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
