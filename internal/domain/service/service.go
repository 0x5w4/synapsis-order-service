package service

import (
	"order-service/config"
	"order-service/internal/adapter/repository"
	"order-service/pkg/logger"
	"order-service/proto/pb"
)

var _ Service = (*service)(nil)

type Service interface {
	Order() OrderService
}

type Properties struct {
	Config                 *config.Config
	Repo                   repository.Repository
	Logger                 logger.Logger
	InventoryServiceClient pb.InventoryServiceClient
}

type service struct {
	Properties
	orderService OrderService
}

func NewService(
	config *config.Config,
	repo repository.Repository,
	logger logger.Logger,
	inventoryServiceClient pb.InventoryServiceClient,
) (*service, error) {
	props := Properties{
		Config:                 config,
		Repo:                   repo,
		Logger:                 logger,
		InventoryServiceClient: inventoryServiceClient,
	}

	return &service{
		Properties:   props,
		orderService: NewOrderService(props),
	}, nil
}

func (s *service) Order() OrderService {
	return s.orderService
}
