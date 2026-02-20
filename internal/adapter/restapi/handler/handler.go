package handler

import (
	"errors"
	"fmt"
	"order-service/config"
	"order-service/internal/domain/service"
	"order-service/internal/shared"
	"order-service/pkg/logger"

	validator "github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type Handler interface {
	Order() OrderHandler
}

type properties struct {
	config    *config.Config
	logger    logger.Logger
	service   service.Service
	validator *validator.Validate
	db        *bun.DB
}

type handler struct {
	properties
	orderHandler OrderHandler
}

func NewHandler(config *config.Config, logger logger.Logger, service service.Service, db *bun.DB) (*handler, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	validate, err := shared.NewValidator()
	if err != nil {
		return nil, fmt.Errorf("failed to setup validator: %w", err)
	}

	props := properties{
		config:    config,
		service:   service,
		logger:    logger,
		validator: validate,
		db:        db,
	}

	h := &handler{
		properties:   props,
		orderHandler: NewOrderHandler(props),
	}

	return h, nil
}

func (h *handler) Order() OrderHandler {
	return h.orderHandler
}
