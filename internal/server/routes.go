package server

import (
	repository "github.com/ilhamgepe/tengahin/internal/repository/postgres"
	authHandlers "github.com/ilhamgepe/tengahin/internal/server/handlers/auth"
	middlewareManager "github.com/ilhamgepe/tengahin/internal/server/middleware"
	"github.com/ilhamgepe/tengahin/internal/service"
	"github.com/ilhamgepe/tengahin/pkg/token"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MountRoutes() {
	v1 := s.echo.Group("/v1")

	// token maker
	tokenMaker := token.NewJwtmaker(s.cfg.Server.JWTSecretKey, s.cfg.Server.JWTRefreshSecretKey)

	// init middleware
	mm := middlewareManager.NewMiddlewareManager(s.logger, s.db, s.rdb, tokenMaker)
	s.echo.Use(middleware.Recover())
	s.echo.Use(mm.Zerolog)

	// init repo
	userRepo := repository.NewUserRepo(s.db)

	// init service
	userService := service.NewUserService(userRepo)
	// init handler
	authHandler := authHandlers.NewAuthHandler(userService, s.rdb, tokenMaker, s.cfg)

	v1GuestWeb := v1.Group("")
	v1WithAuthWeb := v1.Group("", mm.JWTMiddleware)

	// guest
	v1GuestWeb.POST("/register", authHandler.Register)
	v1GuestWeb.POST("/login", authHandler.Login)
	v1GuestWeb.POST("/refresh", authHandler.RefreshToken)
	// with auth
	v1WithAuthWeb.GET("/me", authHandler.Me)
}
