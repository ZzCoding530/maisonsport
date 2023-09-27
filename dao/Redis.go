package dao

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var RedisDB *redis.Client

func InitRedis() {
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Add your Redis password if any
		DB:       0,
	}

	RedisDB = redis.NewClient(options)

	// Ping the Redis server to check if it's running
	pong, err := RedisDB.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
}
