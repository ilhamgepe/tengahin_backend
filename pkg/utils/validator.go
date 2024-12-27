package utils

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/ilhamgepe/validationMessageHelper"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, v interface{}) error {
	return validate.StructCtx(ctx, v)
}

func HandleValidatorError(c echo.Context, err error) error {
	if _, ok := err.(validator.ValidationErrors); ok {
		errors := validationMessageHelper.GenerateMessage(err)
		return c.JSON(http.StatusBadRequest, httpresponse.RestError{
			ErrError:  echo.ErrBadRequest.Error(),
			ErrCauses: errors,
		})
	}
	return c.JSON(http.StatusInternalServerError, httpresponse.RestError{
		ErrError:  echo.ErrInternalServerError.Error(),
		ErrCauses: err.Error(),
	})
}
