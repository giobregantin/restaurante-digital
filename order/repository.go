package order

import "github.com/hsxflowers/restaurante-digital/order/domain"

type OrderRepository struct {
	database domain.OrderStorage
}

func NewOrderRepository(database domain.OrderStorage) *OrderRepository {
	return &OrderRepository{
		database: database,
	}
}
