package model

import (
	"order-service/internal/domain/entity"

	"github.com/uptrace/bun"
)

type Order struct {
    bun.BaseModel `bun:"table:orders,alias:order"`
	Base
    UserID     uint32      `bun:"user_id,notnull"`
    Status     string      `bun:"status,notnull"`
    TotalPrice float64     `bun:"total_price,notnull"`

    Items      []*OrderItem `bun:"rel:has-many,join:id=order_id"`
}

func (m *Order) ToDomain() *entity.Order {
	if m == nil {
		return nil
	}

	return &entity.Order{
		Base: entity.Base{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		UserID:     m.UserID,
		Status:     m.Status,
		TotalPrice: m.TotalPrice,
		Items:      ToOrderItemsDomain(m.Items),
	}
}

func ToOrdersDomain(arg []*Order) []*entity.Order {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*entity.Order, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, arg[i].ToDomain())
	}

	return res
}

func AsOrder(arg *entity.Order) *Order {
	if arg == nil {
		return nil
	}

	return &Order{
		Base: Base{
			ID:        arg.ID,
			CreatedAt: arg.CreatedAt,
			UpdatedAt: arg.UpdatedAt,
			DeletedAt: arg.DeletedAt,
		},
		UserID:     arg.UserID,
		Status:     arg.Status,
		TotalPrice: arg.TotalPrice,
		Items:      AsOrderItems(arg.Items),
	}
}

func AsOrders(arg []*entity.Order) []*Order {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*Order, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, AsOrder(arg[i]))
	}

	return res
}
