package handlers

import (
	"net/http"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/ilhamgepe/tengahin/internal/model"
	"github.com/ilhamgepe/tengahin/internal/service"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/ilhamgepe/tengahin/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type RoleHandler struct {
	roleService service.RoleService
	cfg         *config.Config
	logger      *zerolog.Logger
}

func NewRoleHandler(roleService service.RoleService, cfg *config.Config, logger *zerolog.Logger) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		cfg:         cfg,
		logger:      logger,
	}
}

func (h *RoleHandler) CreateRole(c echo.Context) error {
	h.logger.Info().Msg("create role in handler")
	var req model.CreateRoleDTO
	if err := utils.ReadRequest(c, &req); err != nil {
		h.logger.Info().Msg("error validation in handler")
		return utils.HandleValidatorError(c, err)
	}
	h.logger.Info().Msg("creating role with service in handler")
	err := h.roleService.CreateRole(c.Request().Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("error creating role in handler")
		return httpresponse.KnownSQLError(c, err)
	}
	h.logger.Info().Msg("no error here")
	return c.NoContent(http.StatusOK)
}
