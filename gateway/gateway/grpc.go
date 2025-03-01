package gateway

import (
	"context"
	"log"

	pb "github.com/mrd1920/oms-common/api"
	"github.com/mrd1920/oms-common/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	client, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to connect to any orders service: %v", err)
		return nil, err
	}

	// c := pb.NewOrderServiceClient(client)
	c := pb.NewOrderManagementClient(client)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})

}
