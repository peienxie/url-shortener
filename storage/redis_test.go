package storage

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/peienxie/url-shortener/config"
	"github.com/peienxie/url-shortener/shorten"
	"github.com/stretchr/testify/require"
)

var store URLStore

func init() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("load config app.env err: %v\n", err)
	}

	client := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatalf("connect redis err: %v\n", err)
	}
	store = NewRedisURLStore(client)
}

func TestSaveThenLoadURL(t *testing.T) {
	longURL := "https://amazon.com"
	shortURL := shorten.ShortenByHash(longURL)

	err := store.SaveURL(context.Background(), shortURL, longURL, time.Second)
	require.NoError(t, err)

	loadedURL, err := store.LoadURL(context.Background(), shortURL)
	require.NoError(t, err)
	require.Equal(t, longURL, loadedURL)
}

func TestLoadExpiredURL(t *testing.T) {
	longURL := "https://google.com"
	shortURL := shorten.ShortenByHash(longURL)

	err := store.SaveURL(context.Background(), shortURL, longURL, time.Millisecond)
	require.NoError(t, err)

	time.Sleep(time.Millisecond * 2)

	loadedURL, err := store.LoadURL(context.Background(), shortURL)
	require.Error(t, err)
	require.Equal(t, redis.Nil, err)
	require.GreaterOrEqual(t, "", loadedURL)
}

func TestLoadByWrongShortenURL(t *testing.T) {
	longURL := "https://facebook.com"
	shortURL := shorten.ShortenByHash(longURL)

	err := store.SaveURL(context.Background(), shortURL, longURL, time.Second)
	require.NoError(t, err)

	loadedURL, err := store.LoadURL(context.Background(), "http://notthewebsite.com")
	require.Error(t, err)
	require.Equal(t, redis.Nil, err)
	require.Equal(t, "", loadedURL)
}
