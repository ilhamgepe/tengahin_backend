package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/ilhamgepe/tengahin/internal/service"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/ilhamgepe/tengahin/pkg/utils"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	WebAuthSessionName = "tengahin_session"
	WebAuthSessionID   = "session_id"
	UserCtxKey         = "user"
)

type AuthHandler struct {
	userService service.UserService
	rdb         *redis.Client
}

func NewAuthHandler(userService service.UserService, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		rdb:         rdb,
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

	user, err := h.userService.Register(c.Request().Context(), req)
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

	user, err := h.userService.Login(c.Request().Context(), req.Email)
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

	sess, err := session.Get(WebAuthSessionName, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	id := uuid.New().String()

	sess.Values[WebAuthSessionID] = id

	userData, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: "failed to marshal user data",
		})
	}

	res := h.rdb.Set(c.Request().Context(), id, userData, time.Minute*60)
	if res.Err() != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: res.Err().Error(),
		})
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
			ErrError:  echo.ErrInternalServerError.Error(),
			ErrCauses: err.Error(),
		})
	}

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	user := c.Get(UserCtxKey)

	return c.JSON(200, httpresponse.RestSuccess{
		Status: http.StatusOK,
		Data: map[string]interface{}{
			"user": user,
		},
	})
}
