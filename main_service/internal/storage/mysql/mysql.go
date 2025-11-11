package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"main_service/internal/models"
	"main_service/internal/storage"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func New(dsn string) (*Repository, error) {
	const op = "mysql.New"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{db: db}, nil
}

// * SaveMetadata сохраняет метаданные для текста
// * ttlDays — время жизни записи в днях
func (r *Repository) SaveMetadata(ctx context.Context, hash string, ttlDays int) error {
	const op = "mysql.SaveMetadata"

	now := time.Now().UTC()
	expires := now.AddDate(0, 0, ttlDays)

	query := `INSERT INTO pastes (hash, created_at, expires_at) VALUES (?, ?, ?)`

	if _, err := r.db.ExecContext(ctx, query, hash, now, expires); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// * GetByHash возвращает метаданные для текста по хэшу
func (r *Repository) GetByHash(ctx context.Context, hash string) (*models.Paste, error) {
	const op = "mysql.GetByHash"

	query := `SELECT hash, created_at, expires_at FROM pastes WHERE hash = ?`

	var p models.Paste
	if err := r.db.QueryRowContext(ctx, query, hash).Scan(&p.Hash, &p.CreatedAt, &p.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrTextNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !p.ExpiresAt.IsZero() && p.ExpiresAt.Before(time.Now().UTC()) {
		return nil, storage.ErrTTLIsExpired
	}

	return &p, nil
}

// * GetExpired возвращает хэш для текстов, ttl которых истёк
func (r *Repository) GetExpired(ctx context.Context) ([]string, error) {
	const op = "mysqlRepository.GetExpired"

	query := `SELECT hash FROM pastes WHERE expires_at <= UTC_TIMESTAMP()`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var hashes []string
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		hashes = append(hashes, hash)
	}

	return hashes, nil
}

// * DeleteByHash удаляет метаданные по хэшу
func (r *Repository) DeleteByHash(ctx context.Context, hash string) error {
	const op = "mysqlRepository.DeleteByHash"

	query := `DELETE FROM pastes WHERE hash = ?`

	if _, err := r.db.ExecContext(ctx, query, hash); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// * Close закрывает соединение с базой данных
func (r *Repository) Close() error {
	return r.db.Close()
}
