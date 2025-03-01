package gateway

import "context"

type OrdersGateway interface {
	UpdateOrderAfterPayment(ctx context.Context, orderID, paymentLink string) error
}
