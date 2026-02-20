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

type properties struct {
	config                 *config.Config
	repo                   repository.Repository
	logger                 logger.Logger
	inventoryServiceClient pb.InventoryServiceClient
}

type service struct {
	properties
	orderService OrderService
}

func NewService(
	config *config.Config,
	repo repository.Repository,
	logger logger.Logger,
	inventoryServiceClient pb.InventoryServiceClient,
) (*service, error) {
	props := properties{
		config:                 config,
		repo:                   repo,
		logger:                 logger,
		inventoryServiceClient: inventoryServiceClient,
	}

	return &service{
		properties:   props,
		orderService: NewOrderService(props),
	}, nil
}

func (s *service) Order() OrderService {
	return s.orderService
}
