package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	"em-internship/internal/config"
	"em-internship/internal/handlers"
	"em-internship/internal/repository"
	"em-internship/internal/service"
)

func main() {
	logger, err := config.NewLogger(&config.LoggingConfig{
		Level:       os.Getenv("LOG_LEVEL"),
		Development: os.Getenv("LOG_DEVELOPMENT") == "true",
	})

	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("failed load config", zap.Error(err))
	}

	db, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		logger.Fatal("failed connect to database", zap.Error(err))
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		logger.Fatal("failed ping database", zap.Error(err))
	}

	m, err := migrate.New("file://migrations", cfg.Database.DSN())
	if err != nil {
		logger.Fatal("failed initialize migration", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("failed initialize migration", zap.Error(err))
	}

	logger.Info("migrations applied successfully")

	subRep := repository.NewSubscriptionRepository(db, logger)
	subService := service.NewSubscriptionService(subRep, logger)
	subHandler := handlers.NewSubscriptionHandler(subService, logger)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", subHandler.CreateSubscription)
		r.Get("/", subHandler.ListSubscriptions)
		r.Get("/{id}", subHandler.GetSubscription)
		r.Put("/{id}", subHandler.UpdateSubscription)
		r.Delete("/{id}", subHandler.DeleteSubscription)
		r.Get("/total-cost", subHandler.GetTotalCost)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting server", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("failed shutdown server", zap.Error(err))
	}

	logger.Info("server stopped")
}
