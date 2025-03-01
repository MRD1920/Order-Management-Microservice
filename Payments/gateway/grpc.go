package gateway

import (
	"context"
	"log"

	pb "github.com/mrd1920/oms-common/api"

	"github.com/mrd1920/oms-common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) UpdateOrderAfterPayment(ctx context.Context, orderID, paymentLink string) error {

	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewOrderManagementClient(conn)

	_, err = ordersClient.UpdateOrder(ctx, &pb.Order{
		Id:          orderID,
		Status:      "waiting_for_shipment",
		PaymentLink: paymentLink,
	})

	return err

}
