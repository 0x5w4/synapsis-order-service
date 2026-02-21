package app

import (
	"context"
	"errors"
	"fmt"
	"order-service/config"
	"order-service/internal/adapter/grpcclient"
	"order-service/internal/adapter/repository"
	rest "order-service/internal/adapter/restapi"
	"order-service/internal/domain/service"
	"order-service/pkg/apmtracer"
	"order-service/pkg/bundb"
	"order-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config     *config.Config
	restServer rest.Server
	logger     logger.Logger
	tracer     apmtracer.Tracer
}

func NewApp(config *config.Config, logger logger.Logger) (*App, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &App{
		config: config,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var (
		err error
	)

	// Initialize tracer
	a.tracer, err = apmtracer.NewApmTracer(&apmtracer.Config{
		ServiceName:    a.config.Tracer.ServiceName,
		ServiceVersion: a.config.Tracer.ServiceVersion,
		ServerURL:      a.config.Tracer.ServerURL,
		SecretToken:    a.config.Tracer.SecretToken,
		Environment:    a.config.Tracer.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize repository
	repo, err := repository.NewRepository(a.config, a.logger)
	if err != nil {
		return fmt.Errorf("failed to setup repository: %w", err)
	}

	// Initialize service
	inventorySvcClient, err := grpcclient.NewInventoryServiceClient(a.config)
	if err != nil {
		return fmt.Errorf("failed to create inventory service client: %w", err)
	}

	service, err := service.NewService(a.config, repo, a.logger, inventorySvcClient)
	if err != nil {
		return fmt.Errorf("failed to setup service: %w", err)
	}

	// Initialize and start REST server
	a.restServer, err = rest.NewEchoServer(a.config, a.logger, service, repo)
	if err != nil {
		return fmt.Errorf("failed to setup server: %w", err)
	}

	if err := a.restServer.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	a.logger.Info().Msgf("Server started at %s:%d", a.config.HTTP.Host, a.config.HTTP.Port)

	// Wait for shutdown signal
	<-ctx.Done()
	a.logger.Info().Msg("Shutdown signal received, starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Shutdown REST server
	if err := a.restServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error().Err(err).Msg("Failed to gracefully shutdown REST server")
	} else {
		a.logger.Info().Msg("REST server shut down gracefully")
	}

	// Close repository
	if err := repo.Close(); err != nil {
		a.logger.Error().Err(err).Msg("Failed to gracefully close repository")
	} else {
		a.logger.Info().Msg("Repository closed gracefully")
	}

	a.tracer.Shutdown()

	return nil
}

func (a *App) Migrate(reset bool) error {
	db, err := bundb.NewBunDB(a.config, a.logger)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			a.logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	if reset {
		if err := db.Reset(); err != nil {
			return err
		}
	} else {
		if err := db.Migrate(); err != nil {
			return err
		}
	}

	return nil
}
