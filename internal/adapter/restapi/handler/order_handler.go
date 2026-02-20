package handler

import (
	"net/http"
	"order-service/internal/adapter/restapi/response"
	"order-service/internal/adapter/restapi/serializer"
	"order-service/internal/domain/entity"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Cancel(c echo.Context) error
}

type orderHandler struct {
	properties
}

func NewOrderHandler(props properties) OrderHandler {
	return &orderHandler{properties: props}
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

func (h *orderHandler) Create(c echo.Context) error {
	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	items := make([]*entity.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order := &entity.Order{
		UserID: 1, // TODO: get from auth
		Items:  items,
	}

	createdOrder, err := h.service.Order().Create(c.Request().Context(), order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, serializer.SerializeOrder(createdOrder))
}

func (h *orderHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}

	order, err := h.service.Order().FindByID(c.Request().Context(), uint32(id))
	if err != nil {
		return err
	}

	return response.Success(c, "Order retrieved successfully", serializer.SerializeOrder(order))
}

func (h *orderHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	orders, total, err := h.service.Order().Find(c.Request().Context(), 1, page, perPage) // TODO: get user id from auth
	if err != nil {
		return err
	}

	return response.Paginate(c, "Orders retrieved successfully", serializer.SerializeOrders(orders), response.Pagination{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		TotalPage:  (total + perPage - 1) / perPage,
	})
}

func (h *orderHandler) Cancel(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}

	err = h.service.Order().Cancel(c.Request().Context(), uint32(id))
	if err != nil {
		return err
	}

	return response.Success(c, "Order cancelled successfully", nil)
}
