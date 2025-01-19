package model

import (
	"time"
)

type Role struct {
	ID        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
}

type CreateRoleDTO struct {
	Name string `json:"name" form:"name" validate:"required" db:"name"`
}
