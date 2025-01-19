package server

import (
	"net/http"

	repository "github.com/ilhamgepe/tengahin/internal/repository/postgres"
	authHandlers "github.com/ilhamgepe/tengahin/internal/server/handlers/auth"
	roleHandlers "github.com/ilhamgepe/tengahin/internal/server/handlers/user"
	middlewareManager "github.com/ilhamgepe/tengahin/internal/server/middleware"
	"github.com/ilhamgepe/tengahin/internal/service"
	httpresponse "github.com/ilhamgepe/tengahin/pkg/httpResponse"
	"github.com/ilhamgepe/tengahin/pkg/oauth"
	"github.com/ilhamgepe/tengahin/pkg/token"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MountRoutes() {
	// token maker
	tokenMaker := token.NewJwtmaker(s.cfg.Server.JWTSecretKey, s.cfg.Server.JWTRefreshSecretKey)

	// oauth init
	oauthProviders := oauth.NewOauthProviders(s.cfg)

	// init middleware
	mm := middlewareManager.NewMiddlewareManager(s.logger, s.db, s.rdb, tokenMaker)
	s.echo.Use(middleware.Recover())
	s.echo.Use(mm.Zerolog)
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// init repo
	userRepo := repository.NewUserRepo(s.db, s.logger)
	roleRepo := repository.NewRoleRepo(s.db, s.logger)

	// init service
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)

	// init handler
	authHandler := authHandlers.NewAuthHandler(userService, s.rdb, tokenMaker, s.cfg, oauthProviders, s.logger)
	roleHandler := roleHandlers.NewRoleHandler(roleService, s.cfg, s.logger)
	v1 := s.echo.Group("/v1")

	v1.GET("/ok", func(c echo.Context) error {
		return c.JSON(http.StatusOK, httpresponse.RestSuccess{
			Status: http.StatusOK,
			Data:   "ok",
		})
	})

	// auth routes
	v1AuthRoutes := v1.Group("/auth")
	v1AuthRoutes.POST("/register", authHandler.Register)
	v1AuthRoutes.POST("/login", authHandler.Login)
	v1AuthRoutes.POST("/refresh", authHandler.RefreshToken)
	v1AuthRoutes.POST("/logout", authHandler.Logout)
	v1AuthRoutes.GET("/google", authHandler.GoogleLogin)
	v1AuthRoutes.GET("/google/callback", authHandler.GoogleCallback)
	v1AuthRoutes.GET("/github", authHandler.GithubLogin)
	v1AuthRoutes.GET("/github/callback", authHandler.GithubCallback)

	// role routes
	v1RoleRoutes := v1.Group("/roles")
	v1RoleRoutes.POST("/create", roleHandler.CreateRole)

	// me
	v1.GET("/me", authHandler.Me, mm.JWTMiddleware)
}
