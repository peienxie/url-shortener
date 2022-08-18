package main

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/peienxie/url-shortener/api"
	"github.com/peienxie/url-shortener/config"
	"github.com/peienxie/url-shortener/storage"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("load config app.env err: %v\n", err)
	}
	log.Printf("loaded conifg: %#v\n", config)

	client := redis.NewClient(&redis.Options{Addr: config.RedisAddr})
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatalf("connect redis err: %v\n", err)
	}

	store := storage.NewRedisURLStore(client)
	server := api.NewServer(store)
	if err = server.Serve(config.ApiServerAddr); err != nil {
		log.Fatalf("server err: %v\n", err)
	}
}
