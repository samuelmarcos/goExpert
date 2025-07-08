package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	ListOrder() ([]Order, error)
}
