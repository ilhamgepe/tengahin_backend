package db

import (
	"fmt"
	"time"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/ilhamgepe/tengahin/db/migration"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func NewPostgresDB(cfg *config.Config, logger *zerolog.Logger) *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlPassword,
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlDbname)

	db, err := sqlx.Connect(cfg.Postgres.PgDriver, dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed do connect db")
	}

	db.DB.SetMaxOpenConns(cfg.Postgres.PgMaxConn)
	db.DB.SetConnMaxLifetime(time.Duration(cfg.Postgres.PgMaxConnLifetime) * time.Second)
	db.DB.SetConnMaxIdleTime(cfg.Postgres.PgMaxIdleTime)
	db.DB.SetConnMaxIdleTime(cfg.Postgres.PgMaxIdleTime)

	if err := db.DB.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping db")
	}
	logger.Info().Msg("db connected")

	// run migration
	migration.RunDBMigration("file://db/migration", dsn, logger)

	return db
}
