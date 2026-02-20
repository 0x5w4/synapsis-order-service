package postgresrepository

import (
	"context"
	"database/sql"
	"order-service/config"
	"order-service/internal/adapter/repository/postgres/model"
	"order-service/pkg/bundb"
	"order-service/pkg/logger"

	"github.com/uptrace/bun"
)

var _ PostgresRepository = (*postgresRepository)(nil)

type RepositoryAtomicCallback func(r PostgresRepository) error

type PostgresRepository interface {
	DB() *bun.DB
	Atomic(ctx context.Context, config *config.Config, fn RepositoryAtomicCallback) error
	Close() error
	Order() OrderRepository
}

type properties struct {
	db     bun.IDB
	logger logger.Logger
}

type postgresRepository struct {
	properties
	orderRepository OrderRepository
}

func NewPostgresRepository(config *config.Config, logger logger.Logger) (*postgresRepository, error) {
	db, err := bundb.NewBunDB(config, logger)
	if err != nil {
		return nil, err
	}

	db.DB().RegisterModel(
		(*model.Order)(nil),
	)

	return create(properties{
		db:     db.DB(),
		logger: logger,
	}), nil
}

func (r *postgresRepository) DB() *bun.DB {
	dbInstance, ok := r.db.(*bun.DB)
	if !ok {
		r.logger.Error().Msg("Failed to assert type *bun.DB for the underlying database instance")
		return nil
	}

	return dbInstance
}

func (r *postgresRepository) Close() error {
	return r.DB().Close()
}

func (r *postgresRepository) Atomic(ctx context.Context, config *config.Config, fn RepositoryAtomicCallback) error {
	err := r.db.RunInTx(
		ctx,
		&sql.TxOptions{Isolation: sql.LevelSerializable},
		func(ctx context.Context, tx bun.Tx) error {
			return fn(create(properties{
				db:     tx,
				logger: r.logger,
			}))
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func create(props properties) *postgresRepository {
	return &postgresRepository{
		properties:      props,
		orderRepository: NewOrderRepository(props.db, props.logger),
	}
}

func (r *postgresRepository) Order() OrderRepository {
	return r.orderRepository
}
