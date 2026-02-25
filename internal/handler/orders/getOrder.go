package orders

import (
	"github.com/gorilla/mux"
	"github.com/matsapkov/wb_kafka_service/pkg/json"
	"net/http"

	"github.com/matsapkov/wb_kafka_service/internal/handler/orders/dto"
)

func (h *Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	orderId := vars["order_id"]

	if orderId == "" {
		json.WriteError(w, http.StatusBadRequest, "order_id is required")
		return
	}

	order, err := h.Usecase.GetOrder(ctx, orderId)

	if err != nil {
		code, msg := h.HandlerError(err)
		json.WriteError(w, code, msg)
		return
	}

	outputDto := &dto.OutputDto{
		OrderUID:    order.OrderUID,
		TrackNumber: order.TrackNumber,
		DateCreated: order.DateCreated,
		Payload:     order.Payload,
		UpdatedAt:   order.UpdatedAt,
	}

	err = json.Write(w, http.StatusOK, outputDto)
	if err != nil {
		code, msg := h.HandlerError(err)
		json.WriteError(w, code, msg)
		return
	}
}
