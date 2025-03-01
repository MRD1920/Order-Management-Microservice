package main

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
	"github.com/mrd1920/oms-payments/gateway"
	"github.com/mrd1920/oms-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor: processor, gateway: gateway}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	// connect to payment processor
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		return "", err
	}
	// Update the order with the link
	err = s.gateway.UpdateOrderAfterPayment(ctx, order.Id, link)
	if err != nil {
		return "", err
	}

	return link, nil
}
