package main

import (
	pb "github.com/mrd1920/oms-common/api"
)

type CreateOrderRequest struct {
	Order         *pb.Order `json:"order"`
	RedirectToURL string    `json:"redirectToURL"`
}
