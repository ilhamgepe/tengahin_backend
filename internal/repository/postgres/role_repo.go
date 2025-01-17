package repository

import (
	"context"

	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type RoleRepo interface {
	CreateRole(ctx context.Context, arg model.CreateRoleDTO) error
	FindRoleByID(ctx context.Context, id int64) (*model.Role, error)
	FindRoleByName(ctx context.Context, name string) (*model.Role, error)
}

type roleRepo struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

func NewRoleRepo(db *sqlx.DB, logger *zerolog.Logger) RoleRepo {
	return &roleRepo{
		db:     db,
		logger: logger,
	}
}

const (
	createRoleQuery = `
		INSERT INTO roles
			(name)
		VALUES
			($1);
	`
	findRoleByIDQuery = `
		SELECT
			id,name,created_at
		FROM
			roles
		WHERE
			id = $1
	`
	findRoleByNameQuery = `
		SELECT
			id,name,created_at
		FROM
			roles
		WHERE
			name = $1
	`
)

func (rr *roleRepo) CreateRole(ctx context.Context, arg model.CreateRoleDTO) error {
	rr.logger.Info().Msgf("create role %s", arg.Name)
	_, err := rr.db.ExecContext(ctx, createRoleQuery, arg.Name)
	if err != nil {
		rr.logger.Error().Err(err).Msg("failed to create role")
		return err
	}
	rr.logger.Info().Msg("role created successfully")
	return nil
}

func (rr *roleRepo) FindRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	if err := rr.db.SelectContext(ctx, &role, findRoleByIDQuery, id); err != nil {
		return nil, err
	}

	return &role, nil
}

func (rr *roleRepo) FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := rr.db.SelectContext(ctx, &role, findRoleByNameQuery, name); err != nil {
		return nil, err
	}

	return &role, nil
}
