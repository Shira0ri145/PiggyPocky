package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // ถ้ามี password ใส่ตรงนี้
		DB:       0,
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Cannot connect to Redis: %v", err)
	}

	log.Println("Connected to Redis.")
}
