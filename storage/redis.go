package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisURLStore struct {
	*redis.Client
}

// NewRedisURLStore creates a new URL store with given address
func NewRedisURLStore(client *redis.Client) *RedisURLStore {
	return &RedisURLStore{Client: client}
}

func (s *RedisURLStore) SaveURL(ctx context.Context, shortURL, longURL string, expireTime time.Duration) error {
	err := s.Client.WithContext(ctx).Set(shortURL, longURL, expireTime).Err()
	return err
}

func (s *RedisURLStore) LoadURL(ctx context.Context, shortURL string) (string, error) {
	longURL, err := s.Client.WithContext(ctx).Get(shortURL).Result()
	return longURL, err
}
