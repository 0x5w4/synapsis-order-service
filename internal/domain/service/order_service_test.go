package service_test

import (
	"context"
	"testing"

	"order-service/constant"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/service"
	"order-service/mocks"
	"order-service/proto/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

// setupOrderTest initializes the service with all required mock layers
func setupOrderTest(t *testing.T) (
	service.OrderService,
	*mocks.MockRepository,
	*mocks.MockPostgresRepository,
	*mocks.MockOrderRepository,
	*mocks.MockInventoryServiceClient,
) {
	mRepo := mocks.NewMockRepository(t)
	mPostgres := mocks.NewMockPostgresRepository(t)
	mOrder := mocks.NewMockOrderRepository(t)
	mInventory := mocks.NewMockInventoryServiceClient(t)

	// Link the Repository layers
	mRepo.EXPECT().Postgres().Return(mPostgres).Maybe()
	mPostgres.EXPECT().Order().Return(mOrder).Maybe()

	// Initialize service with properties
	s := service.NewOrderService(service.Properties{
		Repo:                   mRepo,
		InventoryServiceClient: mInventory,
	})

	return s, mRepo, mPostgres, mOrder, mInventory
}

func TestOrderService_Create_Success(t *testing.T) {
	s, _, _, mOrder, mInventory := setupOrderTest(t)
	ctx := context.Background()

	inputOrder := &entity.Order{
		Items: []*entity.OrderItem{
			{Base: entity.Base{ID: 101}, Quantity: 2},
		},
	}

	// 1. Mock gRPC: Get Product info (Note the 3rd arg for variadic opts)
	mInventory.EXPECT().
		GetProduct(ctx, &pb.GetProductRequest{Id: 101}, mock.Anything).
		Return(&pb.Product{
			Id:    101,
			Stock: 10,
			Price: 50.0,
		}, nil)

	// 2. Mock DB: Create Order
	expectedCreated := &entity.Order{Base: entity.Base{ID: 1}, TotalPrice: 100.0}
	mOrder.EXPECT().
		Create(ctx, mock.MatchedBy(func(o *entity.Order) bool {
			return o.TotalPrice == 100.0 && o.Status == string(constant.OrderStatusConfirmed)
		})).
		Return(expectedCreated, nil)

	// Execute
	result, err := s.Create(ctx, inputOrder)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), result.ID)
	assert.Equal(t, 100.0, result.TotalPrice)
}

func TestOrderService_Create_StockShortage(t *testing.T) {
	s, _, _, _, mInventory := setupOrderTest(t)
	ctx := context.Background()

	inputOrder := &entity.Order{
		Items: []*entity.OrderItem{{Base: entity.Base{ID: 101}, Quantity: 5}},
	}

	// Mock gRPC: Return stock less than requested
	mInventory.EXPECT().
		GetProduct(ctx, mock.Anything, mock.Anything).
		Return(&pb.Product{Id: 101, Stock: 2, Price: 50.0}, nil)

	result, err := s.Create(ctx, inputOrder)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "stock is not enough")
}

func TestOrderService_Cancel_Success(t *testing.T) {
	s, _, _, mOrder, mInventory := setupOrderTest(t)
	ctx := context.Background()
	orderID := uint32(1)

	// 1. Mock FindByID (Internal call within Cancel)
	existingOrder := &entity.Order{
		Base:   entity.Base{ID: orderID},
		Status: string(constant.OrderStatusConfirmed),
		Items:  []*entity.OrderItem{{Base: entity.Base{ID: 500}}},
	}
	mOrder.EXPECT().FindByID(ctx, orderID).Return(existingOrder, nil)

	// 2. Mock gRPC: Update Inventory Status
	mInventory.EXPECT().
		UpdateReservationStatus(ctx, mock.Anything, mock.Anything).
		Return(&emptypb.Empty{}, nil)

	// 3. Mock DB: Update Order Status
	mOrder.EXPECT().
		UpdateStatus(ctx, orderID, string(constant.OrderStatusCancelled)).
		Return(nil)

	err := s.Cancel(ctx, orderID)

	assert.NoError(t, err)
}

func TestOrderService_FindByID_NotFound(t *testing.T) {
	s, _, _, mOrder, _ := setupOrderTest(t)
	ctx := context.Background()

	mOrder.EXPECT().FindByID(ctx, uint32(999)).Return(nil, nil)

	result, err := s.FindByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found")
}
