package usecase

import (
	"samuelmarcos/goexpert.com/internal/entity"
	"samuelmarcos/goexpert.com/pkg/events"
)

type ListOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (l *ListOrderUseCase) Execute() ([]ListOrderOutputDTO, error) {
	orders, err := l.OrderRepository.ListOrder()
	if err != nil {
		return nil, err
	}

	// Converter []entity.Order para []ListOrderOutputDTO
	dto := make([]ListOrderOutputDTO, len(orders))
	for i, order := range orders {
		dto[i] = ListOrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
	}

	l.OrderCreated.SetPayload(dto)
	l.EventDispatcher.Dispatch(l.OrderCreated)

	return dto, nil
}
