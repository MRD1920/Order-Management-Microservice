package main

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.Order) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) error
	Get(ctx context.Context, id, customerID string) (*pb.Order, error)
	Update(ctx context.Context, id string, order *pb.Order) error
}
