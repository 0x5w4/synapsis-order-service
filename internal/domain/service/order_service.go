package service

import (
	"context"
	"order-service/constant"
	postgresrepository "order-service/internal/adapter/repository/postgres"
	"order-service/internal/domain/entity"
	"order-service/internal/shared/exception"
	"order-service/proto/pb"
)

var _ OrderService = (*orderService)(nil)

type OrderService interface {
	FindByID(ctx context.Context, id uint32) (*entity.Order, error)
	Find(ctx context.Context, userID uint32, page, perPage int) ([]*entity.Order, int, error)
	Create(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Cancel(ctx context.Context, id uint32) error
}

type orderService struct {
	Properties
}

func NewOrderService(props Properties) *orderService {
	return &orderService{
		Properties: props,
	}
}

func (s *orderService) FindByID(ctx context.Context, id uint32) (*entity.Order, error) {
	order, err := s.Repo.Postgres().Order().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, exception.New(exception.TypeNotFound, "404", "order not found")
	}

	return order, nil
}

func (s *orderService) Find(ctx context.Context, userID uint32, page, perPage int) ([]*entity.Order, int, error) {
	orders, total, err := s.Repo.Postgres().Order().Find(ctx, &postgresrepository.FilterOrderPayload{
		UserID:  userID,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (s *orderService) Create(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	var totalPrice float64
	for i, item := range order.Items {
		product, err := s.InventoryServiceClient.GetProduct(ctx, &pb.GetProductRequest{Id: item.ID})
		if err != nil {
			return nil, err
		}

		if product.GetStock() < int32(item.Quantity) {
			return nil, exception.New(exception.TypeBadRequest, "400", "stock is not enough")
		}

		order.Items[i].Price = product.GetPrice()
		order.Items[i].Subtotal = product.GetPrice() * float64(item.Quantity)
		totalPrice += order.Items[i].Subtotal
	}

	order.TotalPrice = totalPrice
	order.Status = string(constant.OrderStatusConfirmed)

	createdOrder, err := s.Repo.Postgres().Order().Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (s *orderService) Cancel(ctx context.Context, id uint32) error {
	order, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != string(constant.OrderStatusConfirmed) {
		return exception.New(exception.TypeBadRequest, "400", "order cannot be cancelled")
	}

	var reservationIDs []uint32
	for _, item := range order.Items {
		reservationIDs = append(reservationIDs, item.ID)
	}

	_, err = s.InventoryServiceClient.UpdateReservationStatus(ctx, &pb.UpdateReservationStatusRequest{
		Ids:    reservationIDs,
		Status: pb.ReservationStatus_RESERVATION_STATUS_CANCELLED,
	})
	if err != nil {
		return err
	}

	return s.Repo.Postgres().Order().UpdateStatus(ctx, id, string(constant.OrderStatusCancelled))
}
