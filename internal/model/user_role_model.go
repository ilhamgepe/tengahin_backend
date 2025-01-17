package model

type UserRole struct {
	ID     int64 `json:"id" db:"id"`
	UserID int64 `json:"user_id" db:"user_id"`
	RoleID int64 `json:"role_id" db:"role_id"`
}

type CreateUserRoleDTO struct {
	ID     int64 `json:"id" db:"id"`
	UserID int64 `json:"user_id" form:"user_id" db:"user_id" validate:"required,min=1"`
	RoleID int64 `json:"role_id" form:"role_id" db:"role_id" validate:"required,min=1"`
}
