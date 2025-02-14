package stripeProcessor

import (
	"fmt"
	"log"

	common "github.com/mrd1920/oms-common"
	pb "github.com/mrd1920/oms-common/api"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var (
	gatewayAddr = common.EnvString("GATEWAY_ADDR", "http://localhost:8080")
)

type StripeProcessor struct {
}

func NewProcessor() *StripeProcessor {
	return &StripeProcessor{}
}

func (s *StripeProcessor) CreatePaymentLink(order *pb.Order) (string, error) {
	log.Printf("Creating payment link for order %v", order)

	gatewaySuccessURL := fmt.Sprintf("%s/success.html", gatewayAddr)

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
	}

	checkoutSession, err := session.New(params)
	if err != nil {
		return "", err
	}
	return checkoutSession.URL, nil
}
