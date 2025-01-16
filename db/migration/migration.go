package migration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/rs/zerolog"
)

func RunDBMigration(migrationUrl, dbsource string, logger *zerolog.Logger) {
	m, err := migrate.New(migrationUrl, dbsource)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create migrate")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal().Err(err).Msg("failed to migrate db")
	}

	logger.Info().Msg("db migrated successfully")
}
