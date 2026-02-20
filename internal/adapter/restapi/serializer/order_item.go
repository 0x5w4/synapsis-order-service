package serializer

import (
	"order-service/internal/domain/entity"
	"time"
)

type OrderItemResponse struct {
	ID        uint32    `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Subtotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SerializeOrderItem(arg *entity.OrderItem) *OrderItemResponse {
	if arg == nil {
		return nil
	}

	return &OrderItemResponse{
		ID:        arg.ID,
		ProductID: arg.ProductID,
		Quantity:  arg.Quantity,
		Price:     arg.Price,
		Subtotal:  arg.Subtotal,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
	}
}

func SerializeOrderItems(arg []*entity.OrderItem) []*OrderItemResponse {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*OrderItemResponse, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, SerializeOrderItem(arg[i]))
	}

	return res
}
