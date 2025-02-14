package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/mrd1920/oms-common/api"
	"github.com/mrd1920/oms-common/broker"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	channel *amqp091.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp091.Channel) {

	handler := &grpcHandler{service: service, channel: channel}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

// TODO: Implement the CreateOrder method of the OrderServiceServer interface from oms-common/api/oms_grpc.pb.go
func (h *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! Order %v", payload)

	o, err := h.service.CreateOrder(ctx, payload)
	if err != nil {
		return nil, err;
	}
	
	// o := &pb.Order{
	// 	Id: "123",
	// }

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp091.Persistent,
	})
	return o, nil
}
