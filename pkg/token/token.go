package token

import (
	"errors"
	"time"
)

var (
	ErrInvalidTokenMethod = errors.New("unexpected signing method")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("token is invalid")
	ErrParsingToken       = errors.New("failed to parse token")
)

type Maker interface {
	CreateToken(id int64, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
	CreateRefreshToken(id int64, duration time.Duration) (string, *Payload, error)
	VerifyRefreshToken(token string) (*Payload, error)
}
