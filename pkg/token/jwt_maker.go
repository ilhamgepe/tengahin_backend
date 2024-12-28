package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type JWTMaker struct {
	secret        string
	refreshSecret string
}

func NewJwtmaker(secret string, refreshSecret string) Maker {
	return &JWTMaker{secret: secret, refreshSecret: refreshSecret}
}

func (m *JWTMaker) CreateToken(id int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(id, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(m.secret))
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenMethod
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrParsingToken
	}

	if valid := jwtToken.Valid; !valid {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (m *JWTMaker) CreateRefreshToken(id int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(id, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(m.refreshSecret))
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (m *JWTMaker) VerifyRefreshToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenMethod
		}
		return []byte(m.refreshSecret), nil
	})
	log.Error().Err(err).Msg("error")
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrParsingToken
	}

	if valid := jwtToken.Valid; !valid {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
