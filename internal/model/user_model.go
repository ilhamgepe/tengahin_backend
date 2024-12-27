package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int64     `json:"id" db:"id"`
	Email            string    `json:"email" db:"email"`
	Username         string    `json:"username" db:"username"`
	Fullname         string    `json:"fullname" db:"fullname"`
	Password         string    `json:"password,omitempty" db:"password"`
	PasswordChangeAt time.Time `json:"password_change_at,omitempty" db:"password_change_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

func (u *User) Sanitize() *User {
	u.Password = ""
	return u
}

func (u *User) HashPassword() error {
	byte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(byte)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
