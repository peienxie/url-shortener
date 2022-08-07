package main

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/peienxie/url-shortener/api"
	"github.com/peienxie/url-shortener/storage"
)

const (
	RedisAddr     = ":6379"
	ApiServerAddr = ":8080"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: RedisAddr})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("connect redis err: %v\n", err)
	}

	store := storage.NewRedisURLStore(client)
	server := api.NewServer(store)
	if err = server.Serve(ApiServerAddr); err != nil {
		log.Fatalf("server err: %v", err)
	}
}
