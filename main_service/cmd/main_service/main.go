package main

import (
	"context"
	"log/slog"
	"main_service/internal/config"
	textService "main_service/internal/http-server/handlers/middleware/text"
	"main_service/internal/http-server/handlers/text/get"
	"main_service/internal/http-server/handlers/text/save"
	kafkaReader "main_service/internal/kafka"
	cleanup "main_service/internal/scheduler"
	minioStorage "main_service/internal/storage/minio"
	"main_service/internal/storage/mysql"
	"main_service/internal/storage/redis"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")

	log := setupLogger(cfg.Env)

	log.Info("starting main service", slog.String("env", cfg.Env))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Info("Shutdown signal received")
		cancel()
	}()

	db, err := mysql.New(cfg.MySQL.DSN)
	if err != nil {
		log.Error("failed to connect mysql", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	blobStorage, err := minioStorage.New(ctx,
		cfg.MinIO.Endpoint,
		cfg.MinIO.User,
		cfg.MinIO.Password,
		cfg.MinIO.Bucket,
		cfg.MinIO.UseSSL,
	)
	if err != nil {
		log.Error("failed to connect minio", slog.String("err", err.Error()))
		os.Exit(1)
	}

	cache, err := redis.New(ctx, cfg.Redis.Db, cfg.Redis.Addr)
	if err != nil {
		log.Error("failed to connect redis", slog.String("err", err.Error()))
		os.Exit(1)
	}

	reader := kafkaReader.New(cfg.Kafka.Addr, cfg.Kafka.Topic)

	textService := textService.New(db, reader, blobStorage, cache, cfg.Redis.PopularityThreshold)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/save", save.New(log, ctx, textService, cfg.DefaultTTL))
	r.Get("/{hash}", get.New(log, ctx, textService))

	cleaner := cleanup.New(db, blobStorage, cache, log)

	go cleaner.Start(ctx, 3, 0)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Info("HTTP server is running")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed", slog.String("err", err.Error()))
			cancel()
		}
	}()

	<-ctx.Done()

	log.Info("Shutting down HTTP server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", slog.String("err", err.Error()))
	} else {
		log.Info("Server stopped gracefully")
	}

	log.Info("Main service stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
