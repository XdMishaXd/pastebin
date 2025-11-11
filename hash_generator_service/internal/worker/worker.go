package worker

import (
	"context"
	"log/slog"
	"pastebin/internal/hashgen"
	kafkaWriter "pastebin/internal/kafka"
	sl "pastebin/internal/lib/logger"
	"time"
)

type Worker struct {
	ID         int
	HashLength int
	Producer   *kafkaWriter.KafkaWriter
	Rate       int
	BatchSize  int
}

func New(id, hashLen, rate, batchSize int, producer *kafkaWriter.KafkaWriter) *Worker {
	return &Worker{
		ID:         id,
		HashLength: hashLen,
		Producer:   producer,
		BatchSize:  batchSize,
		Rate:       rate,
	}
}

// * Run запускает worker
func (w *Worker) Run(ctx context.Context, log *slog.Logger) {
	gen := hashgen.New(w.ID, w.HashLength)

	interval := time.Second / time.Duration(w.Rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	messages := make([]string, 0, w.BatchSize)

	for {
		select {
		case <-ctx.Done():
			log.Info("worker %d: stopped", slog.Int("id", w.ID))
			return
		case <-ticker.C:
			hash, err := gen.Generate(w.HashLength)
			if err != nil {
				log.Error("failed to generate hash", slog.Int("id", w.ID), sl.Err(err))
				continue
			}

			messages = append(messages, hash)

			if len(messages) >= w.BatchSize {
				err := w.Producer.SendMessages(ctx, w.ID, messages)
				if err != nil {
					log.Error("worker %d: failed to send messages: %v", slog.Int("id", w.ID), sl.Err(err))
				} else {
					log.Info("worker: sent messages",
						slog.Int("id", w.ID),
						slog.Int("amount", len(messages)),
					)
				}
				messages = messages[:0]
			}
		}
	}
}
