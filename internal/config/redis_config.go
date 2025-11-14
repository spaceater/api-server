package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `json:"redis_host"`
	Port     int    `json:"redis_port"`
	Password string `json:"redis_password"`
	DB       int    `json:"redis_db"`
}

var (
	RedisConfigure RedisConfig
	RedisClient    *redis.Client
)

func InitRedisConfig(configData map[string]interface{}) {
	RedisConfigure = RedisConfig{
		Host:     getConfigString(getJSONTag(RedisConfig{}, "Host"), configData, "127.0.0.1"),
		Port:     getConfigInt(getJSONTag(RedisConfig{}, "Port"), configData, 6379),
		Password: getConfigString(getJSONTag(RedisConfig{}, "Password"), configData, ""),
		DB:       getConfigInt(getJSONTag(RedisConfig{}, "DB"), configData, 0),
	}
	initRedisInstance()
}

func initRedisInstance() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", RedisConfigure.Host, RedisConfigure.Port),
		Password: RedisConfigure.Password,
		DB:       RedisConfigure.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Redis connection established successfully: %s:%d", RedisConfigure.Host, RedisConfigure.Port)
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
