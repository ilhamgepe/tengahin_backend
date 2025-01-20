package repository

import (
	"os"
	"testing"
	"time"

	"github.com/ilhamgepe/tengahin/config"
	baseDB "github.com/ilhamgepe/tengahin/db"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	db     *sqlx.DB
	logger zerolog.Logger
	rdb    *redis.Client
)

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../../../config", "config-local")
	if err != nil {
		panic(err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	db = baseDB.NewPostgresDB(cfg, &logger)
	rdb = baseDB.NewRedisClient(cfg, &logger)

	test := m.Run()
	os.Exit(test)
}
