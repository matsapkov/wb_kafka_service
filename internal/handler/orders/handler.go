package orders

import (
	"database/sql"
	"errors"
	"net/http"

	orderUsecase "github.com/matsapkov/wb_kafka_service/internal/usecase/orders"
)

type Handler struct {
	usecase orderUsecase.Usecase
}

func NewOrderHandler(u orderUsecase.Usecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) handlerError(err error) (int, string) {
	switch {
	case err == nil:
		return http.StatusOK, ""
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound, "order not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
