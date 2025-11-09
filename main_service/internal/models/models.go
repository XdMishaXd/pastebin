package models

import (
	"context"
	"time"
)

type Paste struct {
	Hash      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type TextSaver interface {
	SaveText(ctx context.Context, text string, ttl int) (string, error)
	GetText(ctx context.Context, hash string) (string, error)
	DeleteText(ctx context.Context, hash string) error
}
