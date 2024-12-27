package middlewareManager

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (mm *middlewareManager) Zerolog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogMethod: true,
			LogError:  true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				log := mm.logger.Info()

				if v.Status >= 500 {
					log = mm.logger.Error().Err(v.Error)
				} else if v.Status != 200 && v.Status != 201 {
					log = mm.logger.Warn()
				}

				log.Str("URI", v.URI).
					Str("method", v.Method).
					Int("status", v.Status).
					Str("duration", time.Since(v.StartTime).String()).
					Msg("request")

				return nil
			},
		})(next)(c)
	}
}
