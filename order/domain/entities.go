package domain

import (
	"context"
	"time"

	"github.com/hsxflowers/restaurante-digital/exceptions"
)

type OrderStorage interface {
	GetItem(ctx context.Context, tag string) (*Item, error)
	CreateOrder(ctx context.Context, channel *Order) error
}

type OrderDatabase interface {
	CreateOrder(ctx context.Context, cat *Order) error
	GetItem(ctx context.Context, catId string) (*Item, error)
}

type Service interface {
	// GetItem(ctx context.Context, tag string) (*ItemResponse, error)
	CreateOrder(ctx context.Context, cat *OrderRequest) (*OrderResponse, error)
}

type Item struct {
	ItemId            string        `json:"item_id"`
	Nome              string        `json:"nome"`
	TempoCorte        time.Duration `json:"tempo_corte"`
	TempoGrelha       time.Duration `json:"tempo_grelha"`
	TempoMontagem     time.Duration `json:"tempo_montagem"`
	TempoBebida       time.Duration `json:"tempo_bebida"`
	QuantidadeTarefas int           `json:"quantidade_tarefas,omitempty"`
	Cancelamento      chan struct{} `json:"cancelamento,omitempty"`
	TempoEstimado     time.Duration `json:"tempo_estimado,omitempty"`
}

type ItemResponse struct {
	ItemId            string        `json:"item_id"`
	Nome              string        `json:"nome"`
	TempoCorte        time.Duration `json:"tempo_corte"`
	TempoGrelha       time.Duration `json:"tempo_grelha"`
	TempoMontagem     time.Duration `json:"tempo_montagem"`
	TempoBebida       time.Duration `json:"tempo_bebida"`
}

type OrderRequest struct {
	UserId string `json:"user_id"`
	ItemId string `json:"item_id"`
}

type OrderResponse struct {
	UserId        string        `json:"user_id"`
	ItemId        string        `json:"item_id"`
	Status        string        `json:"status"`
	TempoEstimado time.Duration `json:"tempo_estimado,omitempty"`
}

type Order struct {
	UserId string  `json:"user_id"`
	ItemId string  `json:"item_id"`
	Value  float32 `json:"value"`
	Status string  `json:"status"`
}

func (p *OrderRequest) Validate() error {
	if p.UserId == "" {
		return exceptions.ErrUserIdIsRequired
	}

	if p.ItemId == "" {
		return exceptions.ErrItemIdIsRequired
	}

	return nil
}
