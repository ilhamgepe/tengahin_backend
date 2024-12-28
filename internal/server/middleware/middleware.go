package middlewareManager

import (
	"github.com/ilhamgepe/tengahin/pkg/token"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type middlewareManager struct {
	logger     *zerolog.Logger
	db         *sqlx.DB
	rdb        *redis.Client
	tokenMaker token.Maker
}

func NewMiddlewareManager(logger *zerolog.Logger, db *sqlx.DB, rdb *redis.Client, tokenMaker token.Maker) *middlewareManager {
	return &middlewareManager{
		logger:     logger,
		db:         db,
		rdb:        rdb,
		tokenMaker: tokenMaker,
	}
}
