package postgresrepository

import (
	"context"
	"order-service/internal/adapter/repository/postgres/model"
	"order-service/internal/domain/entity"
	"order-service/internal/shared/exception"
	"order-service/pkg/logger"

	"github.com/uptrace/bun"
)

var _ OrderRepository = (*orderRepository)(nil)

type OrderRepository interface {
	FindByID(ctx context.Context, id uint32) (*entity.Order, error)
	Find(ctx context.Context, filter *FilterOrderPayload) ([]*entity.Order, int, error)
	Create(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Delete(ctx context.Context, id uint32) error
	Update(ctx context.Context, order *entity.Order) (*entity.Order, error)
	UpdateStatus(ctx context.Context, id uint32, status string) error
}

type orderRepository struct {
	db     bun.IDB
	logger logger.Logger
}

func NewOrderRepository(db bun.IDB, logger logger.Logger) *orderRepository {
	return &orderRepository{db: db, logger: logger}
}

func (r *orderRepository) GetTableName() string {
	return "orders"
}

type FilterOrderPayload struct {
	IDs     []uint32
	UserID  uint32
	Page    int
	PerPage int
}

func (r *orderRepository) Find(ctx context.Context, filter *FilterOrderPayload) ([]*entity.Order, int, error) {
	var orders []*model.Order

	query := r.db.NewSelect().Model(&orders)

	if len(filter.IDs) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.IDs))
	}

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	totalCount, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, exception.NewDBError(err, r.GetTableName(), "count order")
	}

	if totalCount == 0 {
		return []*entity.Order{}, 0, nil
	}

	if filter.PerPage > 0 {
		query = query.Limit(filter.PerPage)
	}

	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		query = query.Offset(offset)
	}

	query = query.Order("id DESC")
	if err := query.Scan(ctx); err != nil {
		return nil, 0, exception.NewDBError(err, r.GetTableName(), "find order")
	}

	return model.ToOrdersDomain(orders), totalCount, nil
}

func (r *orderRepository) FindByID(ctx context.Context, id uint32) (*entity.Order, error) {
	if id == 0 {
		return nil, exception.ErrIDNull
	}

	order := &model.Order{Base: model.Base{ID: id}}

	if err := r.db.NewSelect().Model(order).WherePK().Scan(ctx); err != nil {
		return nil, exception.NewDBError(err, r.GetTableName(), "find order by id")
	}

	return order.ToDomain(), nil
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if order == nil {
		return nil, exception.ErrDataNull
	}

	dbOrder := model.AsOrder(order)

	_, err := r.db.NewInsert().Model(dbOrder).Exec(ctx)
	if err != nil {
		return nil, exception.NewDBError(err, r.GetTableName(), "create order")
	}

	return dbOrder.ToDomain(), nil
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if order == nil || order.Base.ID == 0 {
		return nil, exception.ErrDataNull
	}

	dbOrder := model.AsOrder(order)

	_, err := r.db.NewUpdate().Model(dbOrder).WherePK().Exec(ctx)
	if err != nil {
		return nil, exception.NewDBError(err, r.GetTableName(), "update order")
	}

	return dbOrder.ToDomain(), nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uint32, status string) error {
	if id == 0 {
		return exception.ErrIDNull
	}

	dbOrder := &model.Order{Base: model.Base{ID: id}}

	_, err := r.db.NewUpdate().
		Model(dbOrder).
		Set("status = ?", status).
		WherePK().
		Exec(ctx)
	if err != nil {
		return exception.NewDBError(err, r.GetTableName(), "update order status")
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id uint32) error {
	if id == 0 {
		return exception.ErrIDNull
	}

	dbOrder := &model.Order{Base: model.Base{ID: id}}

	_, err := r.db.NewDelete().Model(dbOrder).WherePK().Exec(ctx)
	if err != nil {
		return exception.NewDBError(err, r.GetTableName(), "delete order")
	}

	return nil
}
