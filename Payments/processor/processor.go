package processor

import (
	pb "github.com/mrd1920/oms-common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
