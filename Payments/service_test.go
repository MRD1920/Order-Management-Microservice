package main

import (
	"context"
	"testing"

	"github.com/mrd1920/oms-common/api"
	inmemRegistry "github.com/mrd1920/oms-common/discovery/inmem"
	"github.com/mrd1920/oms-payments/gateway"
	"github.com/mrd1920/oms-payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()

	gateway := gateway.NewGateway(registry)
	svc := NewService(processor, gateway)

	t.Run("should create a payment Link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Errorf("CreatePayment() link is empty")
		}
	})
}
