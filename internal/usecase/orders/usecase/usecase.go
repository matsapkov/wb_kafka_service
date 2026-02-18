package usecase

type Usecase struct {
	orderRepo orders.OrderRepository
}

func NewOrderUsecase(orderRepo orders.OrderRepository) *Usecase {
	return &Usecase{
		orderRepo: orderRepo,
	}
}
