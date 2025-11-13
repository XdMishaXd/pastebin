package kafkaReader

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
}

func New(addr, topic string) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers:           []string{addr},
				Topic:             topic,
				StartOffset:       kafka.LastOffset,
				MaxBytes:          10e6,
				CommitInterval:    30 * time.Second,
				SessionTimeout:    45 * time.Second,
				HeartbeatInterval: 5 * time.Second,
				MaxWait:           10 * time.Second,
				ReadBatchTimeout:  30 * time.Second,
			},
		),
	}
}

// * ReadMessage читает хэши из kafka
func (r *KafkaReader) ReadMessage(ctx context.Context) (string, error) {
	const op = "kafka.ReadMessage"

	msg, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Printf("got from kafka: '%s'", string(msg.Value))

	return string(msg.Value), nil
}

func (r *KafkaReader) Close() error {
	return r.reader.Close()
}
