package middlewareManager

import (
	"encoding/json"
	"net/http"

	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/ilhamgepe/tengahin/internal/server/handlers"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (mm *middlewareManager) GuestWeb(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(handlers.WebAuthSessionName, c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
				ErrError:  echo.ErrInternalServerError.Error(),
				ErrCauses: err.Error(),
			})
		}
		sessionId := sess.Values[handlers.WebAuthSessionID]
		if sessionId != nil {
			res := mm.rdb.Get(c.Request().Context(), sessionId.(string))
			if res.Err() != nil {
				return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
					ErrError:  echo.ErrUnauthorized.Error(),
					ErrCauses: res.Err().Error(),
				})
			}

			return c.JSON(http.StatusForbidden, httpresponse.RestError{
				ErrError:  echo.ErrForbidden.Error(),
				ErrCauses: "this resource is only for guest",
			})
		}

		return next(c)
	}
}

func (mm *middlewareManager) WithAuthWeb(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(handlers.WebAuthSessionName, c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
				ErrError:  echo.ErrInternalServerError.Error(),
				ErrCauses: err.Error(),
			})
		}
		sessionId := sess.Values[handlers.WebAuthSessionID]
		if sessionId == nil || sessionId == "" {
			return c.JSON(http.StatusUnauthorized, httpresponse.RestError{
				ErrError:  echo.ErrUnauthorized.Error(),
				ErrCauses: "unauthorized",
			})
		}

		rData := mm.rdb.Get(c.Request().Context(), sessionId.(string))
		if rData.Err() != nil {
			return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
				ErrError:  echo.ErrInternalServerError.Error(),
				ErrCauses: rData.Err().Error(),
			})
		}

		var user model.User
		if err := json.Unmarshal([]byte(rData.Val()), &user); err != nil {
			return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
				ErrError:  echo.ErrInternalServerError.Error(),
				ErrCauses: err.Error(),
			})
		}

		c.Set(handlers.UserCtxKey, user)

		return next(c)
	}
}
