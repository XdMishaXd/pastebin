package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pastebin/internal/config"
	kafkaWriter "pastebin/internal/kafka"
	"pastebin/internal/worker"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"Starting hash generator",
		slog.String("topic", cfg.Kafka.Topic),
		slog.Int("rate", cfg.Hash.HashRate),
		slog.Int("workers", cfg.Hash.Workers),
		slog.Int("hash_len", cfg.Hash.HashLength),
		slog.Int("batch", cfg.Kafka.BatchSize),
	)

	p := kafkaWriter.New(
		cfg.Kafka.Addr,
		cfg.Kafka.Topic,
		cfg.Kafka.BatchSize,
		cfg.Kafka.MaxAttempts,
	)
	defer p.Close()

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("Shutting down...")
		cancel()
	}()

	for i := 0; i < cfg.Workers; i++ {
		w := worker.New(i, cfg.Hash.HashLength, cfg.Hash.HashRate, cfg.Kafka.BatchSize, p)
		go w.Run(ctx, log)
	}

	<-ctx.Done()
	log.Info("Service gracefully stopped")
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
