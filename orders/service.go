package main

import (
	"context"
	"fmt"

	common "github.com/mrd1920/oms-common"
	pb "github.com/mrd1920/oms-common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store: store}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrNoItems
	}
	mergedItems := mergeItemsQuantities(p.Items)
	fmt.Print(mergedItems)

	//TODO: validate with the stock service
	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)
	freqMap := make(map[string]int32)

	for _, item := range items {
		freqMap[item.ID] += item.Quantity
	}

	for key, val := range freqMap {
		merged = append(merged, &pb.ItemsWithQuantity{
			ID:       key,
			Quantity: val,
		})
	}
	return merged
}
