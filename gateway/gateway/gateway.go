package gateway

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
)

type OrderGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}
