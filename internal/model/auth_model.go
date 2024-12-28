package model

import "golang.org/x/crypto/bcrypt"

type RegisterDTO struct {
	Email    string `json:"email" form:"email" validate:"required,email" db:"email"`
	Username string `json:"username" form:"username" validate:"required" db:"username"`
	Fullname string `json:"fullname" form:"fullname" validate:"required,min=3" db:"fullname"`
	Password string `json:"password" form:"password" validate:"required,min=6" db:"password"`
}

type LoginDTO struct {
	Email    string `json:"email" form:"email" validate:"required,email" db:"email"`
	Password string `json:"password" form:"password" validate:"required,min=6" db:"password"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required"`
}

func (u *RegisterDTO) HashPassword() error {
	byte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(byte)
	return nil
}
