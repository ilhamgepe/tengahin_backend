package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/ilhamgepe/tengahin/internal/service"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/ilhamgepe/tengahin/pkg/oauth"
	"github.com/ilhamgepe/tengahin/pkg/token"
	"github.com/ilhamgepe/tengahin/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	userService   service.UserService
	rdb           *redis.Client
	tokenMaker    token.Maker
	cfg           *config.Config
	oauthProvider *oauth.OauthProviders
}

func NewAuthHandler(userService service.UserService, rdb *redis.Client, tokenMaker token.Maker, cfg *config.Config, oauthProviders *oauth.OauthProviders) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		rdb:           rdb,
		tokenMaker:    tokenMaker,
		cfg:           cfg,
		oauthProvider: oauthProviders,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterDTO
	if err := utils.ReadRequest(c, &req); err != nil {
		return utils.HandleValidatorError(c, err)
	}

	if err := req.HashPassword(); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	user, err := h.userService.CreateUser(c.Request().Context(), req)
	if err != nil {
		return httpresponse.KnownSQLError(c, err)
	}

	user.Sanitize()

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginDTO
	if err := utils.ReadRequest(c, &req); err != nil {
		return utils.HandleValidatorError(c, err)
	}

	user, err := h.userService.FindByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return httpresponse.KnownSQLError(c, err)
	}

	if err := user.ComparePassword(req.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "invalid email or password",
		})
	}

	user.Sanitize()

	token, payload, err := h.tokenMaker.CreateToken(user.ID, h.cfg.Server.TokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	refreshToken, payloadRefresh, err := h.tokenMaker.CreateRefreshToken(user.ID, h.cfg.Server.RefreshTokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	expiresIn := time.Until(payload.ExpiresAt.Time)
	rediRes := h.rdb.Set(c.Request().Context(), payload.ID, userJSON, expiresIn)
	if rediRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: rediRes.Err().Error(),
		})
	}

	expiresIn = time.Until(payloadRefresh.ExpiresAt.Time)
	redisRefreshRes := h.rdb.Set(c.Request().Context(), payloadRefresh.ID, refreshToken, expiresIn)
	if redisRefreshRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: redisRefreshRes.Err().Error(),
		})
	}

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"access_token":             token,
			"expires_at":               payload.ExpiresAt,
			"refresh_token":            refreshToken,
			"refresh_token_expires_at": payloadRefresh.ExpiresAt,
			"user":                     user,
		},
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req model.RefreshTokenDTO
	if err := utils.ReadRequest(c, &req); err != nil {
		return utils.HandleValidatorError(c, err)
	}

	// verify refresh token
	payload, err := h.tokenMaker.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "unauthorized",
		})
	}

	// delete old refresh token yang ada di redis
	redisRes, err := h.rdb.Del(c.Request().Context(), payload.ID).Result()
	if err != nil || redisRes == 0 {
		log.Info().Err(err).Msg("failed to delete refresh token")
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "token not found",
		})
	}

	// ambil id user dan buat int64
	subject, err := payload.GetSubject()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "Invalid refresh token",
		})
	}
	id, err := strconv.Atoi(subject)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "Invalid refresh token",
		})
	}

	// cari user berdasrakan id dari sub yang sudah di convert ke int64
	user, err := h.userService.FindByID(c.Request().Context(), int64(id))
	if err != nil {
		return httpresponse.KnownSQLError(c, err)
	}

	token, payload, err := h.tokenMaker.CreateToken(user.ID, h.cfg.Server.TokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	refreshToken, payloadRefresh, err := h.tokenMaker.CreateRefreshToken(user.ID, h.cfg.Server.RefreshTokenDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	// buat user menjadi json agar bisa di simpan di redis
	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	// simpan user di redis berdasarkan token id
	expiresIn := time.Until(payload.ExpiresAt.Time)
	rediRes := h.rdb.Set(c.Request().Context(), payload.ID, userJSON, expiresIn)
	if rediRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: rediRes.Err().Error(),
		})
	}

	// simpan refresh token di redis untuk kebutuhan refresh, jika logout jangan lupa di hapus
	expiresIn = time.Until(payloadRefresh.ExpiresAt.Time)
	redisRefreshRes := h.rdb.Set(c.Request().Context(), payloadRefresh.ID, refreshToken, expiresIn)
	if redisRefreshRes.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: redisRefreshRes.Err().Error(),
		})
	}

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"access_token":             token,
			"expires_at":               payload.ExpiresAt,
			"refresh_token":            refreshToken,
			"refresh_token_expires_at": payloadRefresh.ExpiresAt,
			"user":                     user,
		},
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	var token model.RefreshTokenDTO
	if err := c.Bind(&token); err != nil {
		return utils.HandleValidatorError(c, err)
	}

	payload, err := h.tokenMaker.VerifyRefreshToken(token.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "unauthorized",
		})
	}

	res := h.rdb.Del(c.Request().Context(), payload.ID)
	log.Info().Any("res", res).Msg("res redis")
	if res.Err() != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "unauthorized",
		})
	}

	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) Me(c echo.Context) error {
	user, ok := c.Get(model.UserCtxKey).(model.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
			ErrError:  echo.ErrUnauthorized.Error(),
			ErrCauses: "unauthorized",
		})
	}

	user.Sanitize()

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"user": user,
		},
	})
}
