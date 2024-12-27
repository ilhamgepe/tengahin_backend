package service

import (
	"context"

	"github.com/ilhamgepe/tengahin/internal/model"
	repository "github.com/ilhamgepe/tengahin/internal/repository/postgres"
)

type UserService interface {
	Register(ctx context.Context, arg model.RegisterDTO) (*model.User, error)
	Login(ctx context.Context, email string) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &authService{
		userRepo: userRepo,
	}
}

func (as *authService) Register(ctx context.Context, arg model.RegisterDTO) (*model.User, error) {
	return as.userRepo.CreateUser(ctx, arg)
}

func (as *authService) Login(ctx context.Context, email string) (*model.User, error) {
	return as.userRepo.FindUserByEmail(ctx, email)
}
