package utils

import (
	"github.com/labstack/echo/v4"
)

// ReadRequest need reference type for request eg: &req
func ReadRequest(c echo.Context, request any) error {
	if err := c.Bind(request); err != nil {
		return err
	}

	return ValidateStruct(c.Request().Context(), request)
}
