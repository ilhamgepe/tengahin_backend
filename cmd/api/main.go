package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/ilhamgepe/tengahin/db"
	"github.com/ilhamgepe/tengahin/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig("./config", "config-local")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Logger

	if cfg.Server.Mode == "development" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.TimeFieldFormat = time.RFC3339
		logger.Info().Msg("mode development")
		// logger.Level(zerolog.DebugLevel)
	}
	rdb := db.NewRedisClient(cfg, &logger)
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close redis")
		}
	}()

	db := db.NewPostgresDB(cfg, &logger)
	defer db.Close()

	server := server.NewServer(&logger, rdb, db, cfg)

	if err := server.Start(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
		return
	}
}
