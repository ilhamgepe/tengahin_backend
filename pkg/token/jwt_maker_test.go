package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

var (
	secret        = "b25d58caa43b0368939fcb31b220594e477c02578bfe27316cceade8571c50b1"
	refreshSecret = "b25d58caa43b0368939fcb31b220594e477c02578bfe27316cceade8571c50b2"
)

func TestJWTMaker(t *testing.T) {
	maker := NewJwtmaker(secret, refreshSecret)
	token, payload, err := maker.CreateToken(1, time.Hour)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload2, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.Equal(t, payload, payload2)
	log.Info().Any("payload", payload).Msg("payload")
	log.Info().Any("payload 2", payload2).Msg("payload 2")
}

func TestJWTExpired(t *testing.T) {
	maker := NewJwtmaker(secret, refreshSecret)
	duration := -time.Minute
	// issuedAt := time.Now()
	// expiredAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(2, duration)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
}

func TestInvalidJWTTokenAlgorithm(t *testing.T) {
	payload, err := NewPayload(2, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker := NewJwtmaker(secret, refreshSecret)
	_, err = maker.VerifyToken(token)
	require.Error(t, err)
}

func TestRefreshJWTMaker(t *testing.T) {
	maker := NewJwtmaker(secret, refreshSecret)
	token, payload, err := maker.CreateRefreshToken(1888, time.Hour)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload2, err := maker.VerifyRefreshToken(token)
	require.Nil(t, err)
	require.Equal(t, payload, payload2)
	log.Info().Any("payload", payload).Msg("payload")
	log.Info().Any("payload 2", payload2).Msg("payload 2")
}

func TestRefreshJWTExpired(t *testing.T) {
	maker := NewJwtmaker(secret, refreshSecret)
	duration := -time.Minute
	// issuedAt := time.Now()
	// expiredAt := issuedAt.Add(duration)

	token, _, err := maker.CreateRefreshToken(2, duration)
	require.NoError(t, err)

	payload, err := maker.VerifyRefreshToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
}

func TestInvalidJWTRefreshTokenAlgorithm(t *testing.T) {
	payload, err := NewPayload(2, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker := NewJwtmaker(secret, refreshSecret)
	_, err = maker.VerifyRefreshToken(token)
	require.Error(t, err)
}
