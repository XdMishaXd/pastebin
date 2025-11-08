package kafkaReader

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
}

func New(addr, topic, groupID string) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: []string{addr},
				Topic:   topic,
				GroupID: groupID,
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

	return string(msg.Value), nil
}

func (r *KafkaReader) Close() error {
	return r.reader.Close()
}
