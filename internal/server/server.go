package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilhamgepe/tengahin/config"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	maxHeaderBytes = 1 << 20 // 1 MB
)

type Server struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *sqlx.DB
	echo   *echo.Echo
	rdb    *redis.Client
}

func NewServer(logger *zerolog.Logger, rdb *redis.Client, db *sqlx.DB, cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		rdb:    rdb,
		db:     db,
		echo:   echo.New(),
	}
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	s.MountRoutes()
	go func() {
		s.logger.Info().Msgf("Server started at %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("failed to start Server")
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.CtxDefaultTimeout*time.Second)
	defer cancel()

	s.logger.Info().Msg("Server is shutting down...")

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Error during server shutdown")
		return err
	}

	s.logger.Info().Msg("Server exited properly")
	return nil
}
