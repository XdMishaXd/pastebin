package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

const popKey = "popular_pastes"

func New(ctx context.Context, db int, addr string) (*RedisRepo, error) {
	const op = "storage.redis.New"

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RedisRepo{client: rdb}, nil
}

// * Text возвращает текст, если он есть в redis
func (r *RedisRepo) Text(ctx context.Context, hash string) (string, error) {
	key := hash

	res, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	return res, err
}

func (r *RedisRepo) SaveText(ctx context.Context, hash, text string) error {
	key := hash

	return r.client.Set(ctx, key, text, 0).Err()
}

func (r *RedisRepo) DeleteText(ctx context.Context, hash string) error {
	key := hash

	return r.client.Del(ctx, key).Err()
}

// * IncPopularity увеличивает популярность конкретного hash
func (r *RedisRepo) IncPopularity(ctx context.Context, hash string) (int64, error) {
	res, err := r.client.ZIncrBy(ctx, popKey, 1, hash).Result()

	return int64(res), err
}

// * Top возвращает топ популярности
func (r *RedisRepo) Top(ctx context.Context, limit int) ([]redis.Z, error) {
	return r.client.ZRevRangeWithScores(ctx, popKey, 0, int64(limit)-1).Result()
}

// * Delete удаляет hash
func (r *RedisRepo) Delete(ctx context.Context, hash string) error {
	return r.client.ZRem(ctx, popKey, hash).Err()
}

func (r *RedisRepo) Close() {
	r.client.Close()
}
