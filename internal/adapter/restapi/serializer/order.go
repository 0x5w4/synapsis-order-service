package serializer

import (
	"order-service/internal/domain/entity"
	"time"
)

type OrderResponse struct {
	ID         uint32               `json:"id"`
	UserID     uint32               `json:"user_id"`
	Status     string               `json:"status"`
	TotalPrice float64              `json:"total_price"`
	Items      []*OrderItemResponse `json:"items"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

func SerializeOrder(arg *entity.Order) *OrderResponse {
	if arg == nil {
		return nil
	}

	return &OrderResponse{
		ID:         arg.ID,
		UserID:     arg.UserID,
		Status:     arg.Status,
		TotalPrice: arg.TotalPrice,
		Items:      SerializeOrderItems(arg.Items),
		CreatedAt:  arg.CreatedAt,
		UpdatedAt:  arg.UpdatedAt,
	}
}

func SerializeOrders(arg []*entity.Order) []*OrderResponse {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*OrderResponse, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, SerializeOrder(arg[i]))
	}

	return res
}
