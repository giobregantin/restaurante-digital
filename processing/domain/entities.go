package domain

import "time"

type Pedido struct {
    PedidoId       string
    UsuarioId      string
    ItemId         string
    Cancelamento   chan struct{}
    Nome           string
    TempoCorte     time.Duration
    TempoGrelha    time.Duration
    TempoMontagem  time.Duration
    TempoBebida    time.Duration
    Valor          float64
    Status         string
    QuantidadeTarefas int
    TempoEstimado  time.Duration
    Prioridade     bool
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

type PedidoDetalhado struct {
    Nome  string
    Valor float64
}