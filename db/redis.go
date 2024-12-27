package db

import (
	"context"
	"fmt"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRedisClient(logger *zerolog.Logger, cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.RedisHost, cfg.Redis.RedisPort),
		Password: cfg.Redis.RedisPassword,
		DB:       cfg.Redis.DB,
		Protocol: cfg.Redis.Protocol,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping redis")
	}

	return rdb
}
