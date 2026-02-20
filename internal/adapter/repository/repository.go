package repository

import (
	"order-service/config"
	postgresrepository "order-service/internal/adapter/repository/postgres"
	"order-service/pkg/logger"
)

type Repository interface {
	Postgres() postgresrepository.PostgresRepository
	Close() error
}

type repository struct {
	postgres postgresrepository.PostgresRepository
}

func NewRepository(config *config.Config, logger logger.Logger) (Repository, error) {
	postgresRepo, err := postgresrepository.NewPostgresRepository(config, logger)
	if err != nil {
		return nil, err
	}

	return &repository{
		postgres: postgresRepo,
	}, nil
}

func (r *repository) Postgres() postgresrepository.PostgresRepository {
	return r.postgres
}

func (r *repository) Close() error {
	return r.postgres.Close()
}
