package kafkaWriter

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func New(addr, topic string, batchSize, maxAttemps int) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers:      []string{addr},
				Topic:        topic,
				Balancer:     &kafka.Hash{},
				BatchSize:    batchSize,
				BatchTimeout: 10 * time.Millisecond,
				MaxAttempts:  maxAttemps,
			},
		),
	}
}

// * SendMessages отправляет []string messages от вокркера workerID
func (p *KafkaWriter) SendMessages(ctx context.Context, workerID int, messages []string) error {
	const op = "kafka.SendMessages"

	kafkaMsgs := make([]kafka.Message, 0, len(messages))

	for _, msg := range messages {
		kafkaMsgs = append(kafkaMsgs, kafka.Message{
			Key:   fmt.Appendf(nil, "w%d-%d", workerID, time.Now().UnixNano()),
			Value: []byte(msg),
			Time:  time.Now(),
		})
	}

	err := p.writer.WriteMessages(ctx, kafkaMsgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *KafkaWriter) Close() error {
	return p.writer.Close()
}
