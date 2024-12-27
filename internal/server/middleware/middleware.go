package middlewareManager

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type middlewareManager struct {
	logger *zerolog.Logger
	db     *sqlx.DB
	rdb    *redis.Client
}

func NewMiddlewareManager(logger *zerolog.Logger, db *sqlx.DB, rdb *redis.Client) *middlewareManager {
	return &middlewareManager{
		logger: logger,
		db:     db,
		rdb:    rdb,
	}
}
