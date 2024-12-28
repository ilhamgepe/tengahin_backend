package middlewareManager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ilhamgepe/tengahin/internal/model"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/labstack/echo/v4"
)

func (mm *middlewareManager) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusForbidden, httpresponse.RestError{
				ErrError:  echo.ErrForbidden.Error(),
				ErrCauses: "who the fuck are you!",
			})
		}
		authorization := strings.Split(token, " ")
		if len(authorization) != 2 || authorization[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
				ErrError:  echo.ErrUnauthorized.Error(),
				ErrCauses: "invalid token format",
			})
		}
		payload, err := mm.tokenMaker.VerifyToken(authorization[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
				ErrError:  echo.ErrUnauthorized.Error(),
				ErrCauses: "unauthorized",
			})
		}

		userJSON, err := mm.rdb.Get(c.Request().Context(), payload.ID).Result()
		if err != nil {
			return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
				ErrError:  echo.ErrUnauthorized.Error(),
				ErrCauses: "unauthorized",
			})
		}

		var user model.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			mm.logger.Warn().Err(err).Msg("Failed to parse user JSON")
			return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
				ErrError:  echo.ErrInternalServerError.Error(),
				ErrCauses: "failed to parse user data",
			})
		}

		c.Set(model.UserCtxKey, user)
		return next(c)
	}
}
