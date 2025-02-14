package main

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(context.Context) error
}
