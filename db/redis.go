package db

import (
	"context"
	"fmt"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRedisClient(cfg *config.Config, logger *zerolog.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.RedisHost, cfg.Redis.RedisPort),
		Password:     cfg.Redis.RedisPassword,
		DB:           cfg.Redis.DB,
		Protocol:     cfg.Redis.Protocol,
		PoolSize:     20, // default 10, (jumlah koneksi ke redis)
		MinIdleConns: 5,  // jumlah koneksi idle minumum
		MaxIdleConns: 10, // jumlah koneksi idle maksimum
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping redis")
	}

	logger.Info().Msg("redis connected")

	return rdb
}
