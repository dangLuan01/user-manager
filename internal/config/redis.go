package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dangLuan01/user-manager/internal/utils"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Port     int
	UserName string
	Password string
	DB       int
}

func NewRedisClient() *redis.Client {
	cfg := RedisConfig{
		Addr:     utils.GetEnv("REDIS_ADDR", "localhost"),
		UserName: utils.GetEnv("REDIS_USER", ""),
		Port:     utils.GetIntEnv("REDIS_PORT", 6379),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       utils.GetIntEnv("REDIS_DB", 0),
	}

	client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Username: cfg.UserName,
        Password: cfg.Password,
        DB:       cfg.DB,
		PoolSize: 20,
		MinIdleConns: 5,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
    })

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("❌ Faile to connecting Redis:%s", err)
	}
	
	log.Println("✅ Connected to Redis...")

	return client
}
