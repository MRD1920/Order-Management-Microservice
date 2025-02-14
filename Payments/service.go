package main

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
	"github.com/mrd1920/oms-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	// connect to payment processor
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		return "", err
	}
	// Update the order with the link

	return link, nil
}
