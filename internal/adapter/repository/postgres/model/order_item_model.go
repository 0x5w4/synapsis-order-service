package model

import (
	"order-service/internal/domain/entity"

	"github.com/uptrace/bun"
)

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items"`
	Base
	OrderID   uint32  `bun:"order_id,notnull"`
	ProductID string  `bun:"product_id,notnull"`
	Quantity  int     `bun:"quantity,notnull"`
	Price     float64 `bun:"price,notnull"`
	Subtotal  float64 `bun:"subtotal,notnull"`

	Order *Order `bun:"rel:belongs-to,join:order_id=id"`
}

func (m *OrderItem) ToDomain() *entity.OrderItem {
	if m == nil {
		return nil
	}

	res := &entity.OrderItem{
		Base: entity.Base{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		ProductID: m.ProductID,
		OrderID:   m.OrderID,
		Quantity:  m.Quantity,
		Price:     m.Price,
		Subtotal:  m.Subtotal,
	}

	if m.Order != nil {
		res.Order = m.Order.ToDomain()
	}

	return res
}

func ToOrderItemsDomain(arg []*OrderItem) []*entity.OrderItem {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*entity.OrderItem, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, arg[i].ToDomain())
	}

	return res
}

func AsOrderItem(arg *entity.OrderItem) *OrderItem {
	if arg == nil {
		return nil
	}

	return &OrderItem{
		Base: Base{
			ID:        arg.ID,
			CreatedAt: arg.CreatedAt,
			UpdatedAt: arg.UpdatedAt,
			DeletedAt: arg.DeletedAt,
		},
		ProductID: arg.ProductID,
		OrderID:   arg.OrderID,
		Quantity:  arg.Quantity,
		Price:     arg.Price,
		Subtotal:  arg.Subtotal,
		Order:     AsOrder(arg.Order),
	}
}

func AsOrderItems(arg []*entity.OrderItem) []*OrderItem {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*OrderItem, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, AsOrderItem(arg[i]))
	}

	return res
}
