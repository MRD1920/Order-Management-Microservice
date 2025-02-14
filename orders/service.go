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

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrder(ctx, p);
	if err != nil {
		return nil, err;
	}
	o:= &pb.Order{
		Id: "42",
		CustomerId: p.CustomerID,
		Status: "pending",
		Items: items,
	}
	return o, nil;
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item,error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}
	mergedItems := mergeItemsQuantities(p.Items)
	fmt.Print(mergedItems)

	//TODO: validate with the stock service

	//Temporary
	var itemsWithPrice []*pb.Item
	for _, item := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1JZ2Z3J2Z3J2Z3J2Z3J2Z3J2",
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	return itemsWithPrice, nil
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
