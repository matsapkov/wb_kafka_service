package orders

type Handler struct {
	Usecase orders.Usecase
}

func NewOrderHandler(u orders.Usecase) *Handler {
	return &Handler{
		Usecase: u,
	}
}

func (h *Handler) HandlerError(err error) (int, string) {

}
