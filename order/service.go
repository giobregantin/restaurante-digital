package order

import (
	"context"

	"github.com/hsxflowers/restaurante-digital/order/domain"
)

type Service struct {
}

func NewOrderService() domain.Service {
	return &Service{}
}

func (s *Service) CreateOrder(ctx context.Context, order *domain.OrderRequest) (*domain.OrderResponse, error) {
	
}
