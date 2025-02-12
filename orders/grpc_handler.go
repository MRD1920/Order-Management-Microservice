package main

import (
	"context"
	"log"

	pb "github.com/mrd1920/oms-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService) {

	handler := &grpcHandler{service: service}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

// TODO: Implement the CreateOrder method of the OrderServiceServer interface from oms-common/api/oms_grpc.pb.go
func (h *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! Order %v", payload)
	o := &pb.Order{
		Id: "123",
	}
	return o, nil
}
