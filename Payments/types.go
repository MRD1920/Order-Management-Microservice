package main

import (
	"context"

	pb "github.com/mrd1920/oms-common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
