package rest

import (
	"context"
	"fmt"
	"net/http"
	"order-service/config"
	"order-service/internal/adapter/repository"
	"order-service/internal/adapter/restapi/handler"
	"order-service/internal/domain/service"
	"order-service/pkg/logger"
	"time"

	"github.com/cockroachdb/errors"
	echo "github.com/labstack/echo/v4"
)

const shutdownTimeout = 10 * time.Second

type Server interface {
	Echo() *echo.Echo
	Start() error
	Shutdown(ctx context.Context) error
}

type echoServer struct {
	config  *config.Config
	logger  logger.Logger
	echo    *echo.Echo
	handler handler.Handler
}

func NewEchoServer(config *config.Config, logger logger.Logger, service service.Service, repository repository.Repository) (*echoServer, error) {
	e := echo.New()
	e.HideBanner = true

	handler, err := handler.NewHandler(config, logger, service, repository.Postgres().DB())
	if err != nil {
		return nil, err
	}

	server := &echoServer{
		config:  config,
		logger:  logger.NewInstance().Field("component", "http_server").Logger(),
		echo:    e,
		handler: handler,
	}

	server.setupMiddlewares()
	server.setupRouter()

	return server, nil
}

func (s *echoServer) Echo() *echo.Echo {
	return s.echo
}

func (s *echoServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.HTTP.Host, s.config.HTTP.Port)
	startErrChan := make(chan error, 1)

	go func() {
		if err := s.echo.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			startErrChan <- errors.Wrapf(err, "server failed to start listening on %s", address)
		} else {
			startErrChan <- nil
		}

		close(startErrChan)
	}()

	select {
	case err := <-startErrChan:
		if err != nil {
			return err
		}

		s.logger.Info().Msg("Server listening")
	case <-time.After(100 * time.Millisecond):
		s.logger.Info().Msg("Server assumed started successfully (listening)")
	}

	return nil
}

func (s *echoServer) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		return errors.Wrap(err, "server shutdown failed")
	}

	return nil
}
