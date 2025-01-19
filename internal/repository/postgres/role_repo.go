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

	CreateUserRole(ctx context.Context, arg model.CreateUserRoleDTO) error
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

func (rr *roleRepo) CreateRole(ctx context.Context, arg model.CreateRoleDTO) error {
	rr.logger.Info().Msgf("create role %s", arg.Name)
	const createRoleQuery = `
		INSERT INTO roles
			(name)
		VALUES
			($1);
	`
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
	const findRoleByIDQuery = `
		SELECT
			id,name,created_at
		FROM
			roles
		WHERE
			id = $1
	`
	if err := rr.db.SelectContext(ctx, &role, findRoleByIDQuery, id); err != nil {
		return nil, err
	}

	return &role, nil
}

func (rr *roleRepo) FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	const findRoleByNameQuery = `
		SELECT
			id,name,created_at
		FROM
			roles
		WHERE
			name = $1
	`
	if err := rr.db.SelectContext(ctx, &role, findRoleByNameQuery, name); err != nil {
		return nil, err
	}

	return &role, nil
}

// !nanti ajadah implemen nya lewat db aja dulu
func (rr *roleRepo) CreateUserRole(ctx context.Context, arg model.CreateUserRoleDTO) error {
	const createUserRoleQuery = `
		INSERT INTO user_roles
			(user_id,role_id)
		VALUES
			($1,$2);
	`
	_, err := rr.db.ExecContext(ctx, createUserRoleQuery, arg.UserID, arg.RoleID)
	if err != nil {
		return err
	}
	return nil
}
