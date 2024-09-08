package domain

import "time"

type Pedido struct {
	PedidoId          string
	UsuarioId         string
	ItemId            string
	Nome              string
	TempoCorte        time.Duration
	TempoGrelha       time.Duration
	TempoMontagem     time.Duration
	TempoBebida       time.Duration
	QuantidadeTarefas int
	Cancelamento      chan struct{}
	TempoEstimado     time.Duration
	Status            string
	Valor             float64
}

type Item struct {
	ItemId        string
	Nome          string
	TempoCorte    time.Duration
	TempoGrelha   time.Duration
	TempoMontagem time.Duration
	TempoBebida   time.Duration
	Valor         float64
}
