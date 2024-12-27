package repository

import (
	"context"
	"time"

	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	CreateUser(ctx context.Context, arg model.RegisterDTO) (*model.User, error)
	FindUserByID(ctx context.Context, id int64) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) CreateUser(ctx context.Context, arg model.RegisterDTO) (*model.User, error) {
	var user model.User

	if err := ur.db.QueryRowxContext(
		ctx, createUserQuery,
		arg.Email,
		arg.Username,
		arg.Fullname,
		arg.Password,
		time.Now(),
	).StructScan(&user); err != nil {
		return nil, err
	}
	// if err := ur.db.QueryRowxContext(
	// 	ctx, createUserQuery,
	// 	arg.Email,
	// 	arg.Username,
	// 	arg.Fullname,
	// 	arg.Password,
	// 	time.Now(),
	// ).StructScan(&user); err != nil {
	// 	return nil, err
	// }

	return &user, nil
}

func (ur *userRepo) FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := ur.db.QueryRowxContext(ctx, findUserByIDQuery, id).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := ur.db.QueryRowxContext(ctx, findUserByEmailQuery, email).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
