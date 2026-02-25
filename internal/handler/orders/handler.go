package orders

import (
	"database/sql"
	"errors"
	"github.com/matsapkov/wb_kafka_service/internal/usecase/orders"
	"net/http"
)

type Handler struct {
	Usecase orders.Usecase
}

func NewOrderHandler(u orders.Usecase) *Handler {
	return &Handler{
		Usecase: u,
	}
}

func (h *Handler) HandlerError(err error) (int, string) {
	switch {
	case err == nil:
		return http.StatusOK, ""
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound, "order not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
