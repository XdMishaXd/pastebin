package cleanup

import (
	"context"
	"log/slog"
	"time"
)

type Storage interface {
	GetExpired(ctx context.Context) ([]string, error)
	DeleteByHash(ctx context.Context, hash string) error
}

type FileStorage interface {
	DeleteFile(ctx context.Context, hash string) error
}

type Cache interface {
	DeleteText(ctx context.Context, hash string) error
}

type Cleaner struct {
	db    Storage
	files FileStorage
	cache Cache
	log   *slog.Logger
}

func New(db Storage, files FileStorage, cache Cache, log *slog.Logger) *Cleaner {
	return &Cleaner{
		db:    db,
		files: files,
		cache: cache,
		log:   log,
	}
}

// * Start запускает проверку в указанное время суток.
func (c *Cleaner) Start(ctx context.Context, hour, minute int) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Until(next)):
				c.run(ctx)
			}
		}
	}()
}

// * run выполняет очистку.
func (c *Cleaner) run(ctx context.Context) {
	c.log.Info("Starting cleanup task...")

	expired, err := c.db.GetExpired(ctx)
	if err != nil {
		c.log.Error("Failed to get expired hashes", slog.Any("error", err))
		return
	}

	if len(expired) == 0 {
		c.log.Info("No expired entries found")
		return
	}

	for _, hash := range expired {
		if err := c.cache.DeleteText(ctx, hash); err != nil {
			c.log.Error("Failed to delete from Redis", slog.String("hash", hash), slog.Any("error", err))
		}

		if err := c.db.DeleteByHash(ctx, hash); err != nil {
			c.log.Error("Failed to delete from MySQL", slog.String("hash", hash), slog.Any("error", err))
			continue
		}

		if err := c.files.DeleteFile(ctx, hash); err != nil {
			c.log.Error("Failed to delete from MinIO", slog.String("hash", hash), slog.Any("error", err))
			continue
		}

		c.log.Info("Deleted expired paste", slog.String("hash", hash))
	}

	c.log.Info("Cleanup task completed", slog.Int("deleted", len(expired)))
}
