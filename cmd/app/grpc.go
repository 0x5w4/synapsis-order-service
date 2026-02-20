package app

import (
	"fmt"
	"order-service/config"
	"order-service/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newInventoryServiceClient(cfg *config.Config) (pb.InventoryServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.GRPC.InventoryHost, cfg.GRPC.InventoryPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	return pb.NewInventoryServiceClient(conn), nil
}
