package orders

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matsapkov/wb_kafka_service/internal/handler/orders/dto"
	"github.com/matsapkov/wb_kafka_service/pkg/json"
)

func (h *Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := mux.Vars(r)["order_id"]
	if orderID == "" {
		json.WriteError(w, http.StatusBadRequest, "order_id is required")
		return
	}

	order, err := h.usecase.GetOrder(ctx, orderID)
	if err != nil {
		code, msg := h.handlerError(err)
		json.WriteError(w, code, msg)
		return
	}

	output := dto.OutputDto{
		OrderUID:    order.OrderUID,
		TrackNumber: order.TrackNumber,
		DateCreated: order.DateCreated,
		Payload:     order.Payload,
		UpdatedAt:   order.UpdatedAt,
	}

	if err := json.Write(w, http.StatusOK, output); err != nil {
		json.WriteError(w, http.StatusInternalServerError, "failed to serialize response")
	}
}
