package storage

import (
	"context"
	"time"
)

type URLStore interface {
	SaveURL(ctx context.Context, shortURL, longURL string, expireTime time.Duration) error
	LoadURL(ctx context.Context, shortURL string) (string, error)
}
