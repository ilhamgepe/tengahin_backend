package service

import (
	"context"

	"github.com/ilhamgepe/tengahin/internal/model"
	repository "github.com/ilhamgepe/tengahin/internal/repository/postgres"
)

type UserService interface {
	CreateUser(ctx context.Context, arg model.RegisterDTO) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx context.Context, arg model.RegisterDTO) (*model.User, error) {
	return u.userRepo.CreateUser(ctx, arg)
}

func (u *userService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.userRepo.FindUserByEmail(ctx, email)
}

func (u *userService) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return u.userRepo.FindUserByID(ctx, id)
}
