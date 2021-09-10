package db

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func ConnectClient() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})
	ok, err := Client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("from db: ", ok)
}
