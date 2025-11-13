package textService

import (
	"context"
	"errors"
	"main_service/internal/models"
	"main_service/internal/storage"
)

type MySql interface {
	SaveMetadata(ctx context.Context, hash string, ttlDays int) error
	GetByHash(ctx context.Context, hash string) (*models.Paste, error)
	GetExpired(ctx context.Context) ([]string, error)
	DeleteByHash(ctx context.Context, hash string) error
}

type Kafka interface {
	ReadMessage(ctx context.Context) (string, error)
}

type MinIO interface {
	SaveStringAsFile(ctx context.Context, hash, content string) error
	GetString(ctx context.Context, hash string) (string, error)
	DeleteFile(ctx context.Context, hash string) error
	ListFiles(ctx context.Context) ([]string, error)
}

type TextSaverService struct {
	mysql MySql
	kafka Kafka
	minio MinIO
}

func New(mysql MySql, k Kafka, min MinIO) *TextSaverService {
	return &TextSaverService{
		mysql: mysql,
		kafka: k,
		minio: min,
	}
}

func (s *TextSaverService) SaveText(ctx context.Context, text string, ttl int) (string, error) {
	hash, err := s.kafka.ReadMessage(ctx)
	if err != nil {
		return "", err
	}

	err = s.minio.SaveStringAsFile(ctx, hash, text)
	if err != nil {
		return "", err
	}

	return hash, s.mysql.SaveMetadata(ctx, hash, ttl)
}

func (s *TextSaverService) GetText(ctx context.Context, hash string) (string, error) {
	_, err := s.mysql.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, storage.ErrTextNotFound) {
			return "", storage.ErrTextNotFound
		}
		if errors.Is(err, storage.ErrTTLIsExpired) {
			return "", storage.ErrTTLIsExpired
		}

		return "", err
	}

	text, err := s.minio.GetString(ctx, hash)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (s *TextSaverService) DeleteText(ctx context.Context, hash string) error {
	if err := s.mysql.DeleteByHash(ctx, hash); err != nil {
		return err
	}

	if err := s.minio.DeleteFile(ctx, hash); err != nil {
		return err
	}

	return nil
}
