package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	_ "song-lib/docs"
	"song-lib/internal/config"
	"song-lib/internal/database/postgres"
	"song-lib/internal/lib/logs"
	"song-lib/internal/services"
	"song-lib/internal/transport/rest/handlers/add"
	"song-lib/internal/transport/rest/handlers/del"
	"song-lib/internal/transport/rest/handlers/get"
	"song-lib/internal/transport/rest/handlers/text"
	"song-lib/internal/transport/rest/handlers/up"
)

func Run() {
	cfg := config.MustLoad()

	log := logs.InitLogger(cfg.LogLevel)
	log.Info("Application starting")

	db, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}
	log.Info("Database connected")
	log.Info("Migration is up")

	src := services.New(db)
	log.Info("Services created")

	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/songs", get.New(log, src))
	router.Get("/songs/{id}/text", text.New(log, src))
	router.Post("/songs", add.New(log, src))
	router.Delete("/songs/{id}", del.New(log, src))
	router.Put("/songs/{id}", up.New(log, src))

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info("starting server", slog.String("address", address))

	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
