package main

import (
	"errors"
	"fmt"
	"net/http"

	common "github.com/mrd1920/oms-common"
	pb "github.com/mrd1920/oms-common/api"
	"github.com/mrd1920/oms-gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	//gateway
	gateway gateway.OrderGateway
}

func NewHandler(gateway gateway.OrderGateway) *handler {
	return &handler{
		gateway: gateway,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() == codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := &CreateOrderRequest{
		Order:         order,
		RedirectToURL: fmt.Sprintf("http://localhost:8080/success.html?customerID=%s&orderID=%s", customerID, order.Id),
	}

	common.WriteJSON(w, http.StatusCreated, res)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, item := range items {
		if item.ID == "" {
			return errors.New("item ID is required")
		}

		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than 0")
		}
	}

	return nil
}
