package orders

import (
	"net/http"

	"github.com/matsapkov/wb_kafka_service/internal/handler/orders/dto"
)

func (h *Handler) GetOrderHandler(r *http.Request, w http.ResponseWriter) {
	ctx := r.Context()

	querryParams := r.URL.Query()
	orderId := querryParams.Get("order_uid")

	order, err := h.Usecase.GetOrder(ctx, orderId)

	if err != nil {

	}

	outputDto := &dto.OutputDto{}

}
