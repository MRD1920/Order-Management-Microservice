package main

import (
	"context"
	"log"

	pb "github.com/mrd1920/oms-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

// TODO: Implement the CreateOrder method of the OrderServiceServer interface from oms-common/api/oms_grpc.pb.go
func (h *grpcHandler) CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("New order received!")
	o := &pb.Order{
		Id: "123",
	}
	return o, nil
}
