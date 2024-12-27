package server

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
	repository "github.com/ilhamgepe/tengahin/internal/repository/postgres"
	"github.com/ilhamgepe/tengahin/internal/server/handlers"
	middlewareManager "github.com/ilhamgepe/tengahin/internal/server/middleware"
	"github.com/ilhamgepe/tengahin/internal/service"
	"github.com/labstack/echo-contrib/session"
)

func init() {
	gob.Register(map[string]interface{}{})
}

func (s *Server) MountRoutes() {
	v1 := s.echo.Group("/v1")

	// init middleware
	mm := middlewareManager.NewMiddlewareManager(s.logger, s.db, s.rdb)
	// s.echo.Use(middleware.Recover())
	s.echo.Use(mm.Zerolog)
	s.echo.Use(session.Middleware(sessions.NewCookieStore([]byte(s.cfg.Server.SessionSecretKey))))

	// init repo
	userRepo := repository.NewUserRepo(s.db)

	// init service
	userService := service.NewUserService(userRepo)

	// init handler
	authHandler := handlers.NewAuthHandler(userService, s.rdb)

	v1GuestWeb := v1.Group("", mm.GuestWeb)
	v1WithAuthWeb := v1.Group("", mm.WithAuthWeb)

	// guest
	v1GuestWeb.POST("/register", authHandler.Register)
	v1GuestWeb.POST("/login", authHandler.Login)

	// with auth
	v1WithAuthWeb.GET("/me", authHandler.Me)
}
