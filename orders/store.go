package main

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
)

type store struct {
	// Add monogDB instance here.
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) error {
	return nil
}
func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {
	order := &pb.Order{}
	return order, nil
}

func (s *store) Update(ctx context.Context, id string, order *pb.Order) error {
	return nil
}
