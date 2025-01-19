package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type UserRepo interface {
	CreateUser(ctx context.Context, arg model.RegisterDTO) (*model.User, error)
	FindUserByID(ctx context.Context, id int64) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepo struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

func NewUserRepo(db *sqlx.DB, logger *zerolog.Logger) UserRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}

func (ur *userRepo) CreateUser(ctx context.Context, arg model.RegisterDTO) (*model.User, error) {
	var user model.User
	var roles model.Role
	const createUserQuery = `
	INSERT INTO users 
		(email,username,fullname,password,password_change_at)
	VALUES
		($1,$2,$3,$4,$5)
	RETURNING id,email,username,fullname,password,password_change_at,created_at;
	`
	const createUserRole = `
	INSERT INTO user_roles
		(user_id,role_id)
	VALUES
		($1,$2);
	`
	tx, err := ur.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: failed to create tx: %w", err)
	}

	// ! create user
	if err := tx.QueryRowxContext(
		ctx, createUserQuery,
		arg.Email,
		arg.Username,
		arg.Fullname,
		arg.Password,
		time.Now(),
	).StructScan(&user); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create user: failed to create user: %w", err)
	}

	// ! create default user role
	// get buyer role
	if err := tx.GetContext(ctx, &roles, "SELECT * FROM roles WHERE name = 'buyer' LIMIT 1"); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create user: failed to find default role: %w", err)
	}

	// ! create user role
	_, err = tx.ExecContext(ctx, createUserRole, user.ID, roles.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create user: failed to create user role: %w", err)
	}

	// ! commit tx
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create user: failed to commit tx: %w", err)
	}

	user.Roles = append(user.Roles, roles)

	return &user, nil
}

func (ur *userRepo) FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	user.Roles = []model.Role{}
	const findUserByIDQuery = `
	SELECT 
		u.id, 
		u.email, 
		u.username, 
		u.fullname, 
		u.password, 
		u.password_change_at, 
		u.created_at,
		r.id AS role_id,
		r.name AS role_name
	FROM 
		users AS u
	JOIN 
		user_roles AS ur ON u.id = ur.user_id
	JOIN 
		roles AS r ON ur.role_id = r.id
	WHERE 
		u.id = $1;
	`
	rows, err := ur.db.QueryxContext(ctx, findUserByIDQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ur.logger.Info().Any("rows!!!!!", rows)
	for rows.Next() {
		roles := struct {
			ID   sql.NullInt64
			Name sql.NullString
		}{}

		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Fullname,
			&user.Password,
			&user.PasswordChangeAt,
			&user.CreatedAt,
			&roles.ID,
			&roles.Name,
		); err != nil {
			ur.logger.Error().Err(err).Msg("Error scanning row")
			return nil, err
		}

		if roles.Name.Valid && roles.ID.Valid {
			user.Roles = append(user.Roles, model.Role{
				ID:   roles.ID.Int64,
				Name: roles.Name.String,
			})
		}
	}
	return &user, nil
}

func (ur *userRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	const findUserByEmailQuery = `
	SELECT 
		u.id, 
		u.email, 
		u.username, 
		u.fullname, 
		u.password, 
		u.password_change_at, 
		u.created_at,
		r.id AS role_id,
		r.name AS role_name
	FROM 
		users AS u
	JOIN 
		user_roles AS ur ON u.id = ur.user_id
	JOIN 
		roles AS r ON ur.role_id = r.id
	WHERE 
		u.email = $1;
	`
	rows, err := ur.db.QueryxContext(ctx, findUserByEmailQuery, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		roles := struct {
			ID   sql.NullInt64
			Name sql.NullString
		}{}

		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.Fullname,
			&user.Password,
			&user.PasswordChangeAt,
			&user.CreatedAt,
			&roles.ID,
			&roles.Name,
		); err != nil {
			ur.logger.Error().Err(err).Msg("Error scanning row")
			return nil, err
		}
		ur.logger.Info().Any("roles", roles).Msg("roles")
		if roles.Name.Valid && roles.ID.Valid {
			user.Roles = append(user.Roles, model.Role{
				ID:   roles.ID.Int64,
				Name: roles.Name.String,
			})
		}
	}

	return &user, nil
}
