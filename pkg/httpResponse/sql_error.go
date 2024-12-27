package httpresponse

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func KnownSQLError(c echo.Context, err error) error {
	log.Error().Err(err).Msg("memek loding")
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, RestError{
			ErrError:  echo.ErrNotFound.Error(),
			ErrCauses: "requested data not found",
		})
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return c.JSON(http.StatusConflict, RestError{
				ErrError:  "duplicate entry",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.ForeignKeyViolation:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "foreign key violation",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.NotNullViolation:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "missing required field",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.CheckViolation:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "check constraint violated",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.StringDataRightTruncationDataException:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "string data right truncation",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.NumericValueOutOfRange:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "numeric value out of range",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.SyntaxError:
			return c.JSON(http.StatusInternalServerError, RestError{
				ErrError:  "syntax error in query",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.SerializationFailure:
			return c.JSON(http.StatusConflict, RestError{
				ErrError:  "transaction conflict or deadlock",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.UndefinedColumn:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "undefined column in query",
				ErrCauses: pgErr.Detail,
			})
		case pgerrcode.UndefinedTable:
			return c.JSON(http.StatusBadRequest, RestError{
				ErrError:  "undefined table in query",
				ErrCauses: pgErr.Detail,
			})
		}
	}

	return c.JSON(http.StatusInternalServerError, RestError{
		ErrError:  "internal server error",
		ErrCauses: pgErr.Message,
	})
}
