package service

import (
	"context"

	"github.com/ilhamgepe/tengahin/internal/model"
	repository "github.com/ilhamgepe/tengahin/internal/repository/postgres"
)

type RoleService interface {
	CreateRole(ctx context.Context, arg model.CreateRoleDTO) error
	FindRoleByID(ctx context.Context, id int64) (*model.Role, error)
	FindRoleByName(ctx context.Context, name string) (*model.Role, error)
}

type roleService struct {
	roleRepo repository.RoleRepo
}

func NewRoleService(roleRepo repository.RoleRepo) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

func (rs *roleService) CreateRole(ctx context.Context, arg model.CreateRoleDTO) error {
	return rs.roleRepo.CreateRole(ctx, arg)
}

func (rs *roleService) FindRoleByID(ctx context.Context, id int64) (*model.Role, error) {
	return rs.roleRepo.FindRoleByID(ctx, id)
}

func (rs *roleService) FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	return rs.roleRepo.FindRoleByName(ctx, name)
}
