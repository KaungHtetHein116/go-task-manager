package redisdb

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	_ = godotenv.Load()

	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("Connected to Redis")
}

func Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return Rdb.Set(Ctx, key, jsonData, 10*time.Minute).Err()
}

// Get
func Get(key string, dest interface{}) (bool, error) {
	val, err := Rdb.Get(Ctx, key).Result()

	if err != nil {
		return false, nil
	}

	return true, json.Unmarshal([]byte(val), dest)
}

func Del(key string) error {
	return Rdb.Del(Ctx, key).Err()
}
