package workers

import (
	"fmt"
	"sync"
	"time"

	"github.com/hsxflowers/restaurante-digital/workers"
)

var Menu = []workers.Pedido{
	{
		Nome:          "Callback Burguer",
		TempoCorte:    3 * time.Second,
		TempoGrelha:   8 * time.Second,
		TempoMontagem: 2 * time.Second,
		TempoBebida:   0,
		Cancelamento:  make(chan struct{}), // Canal de cancelamento inicializado
	},
	{
		Nome:          "Null-Burguer",
		TempoCorte:    4 * time.Second,
		TempoGrelha:   7 * time.Second,
		TempoMontagem: 2 * time.Second,
		TempoBebida:   0,
		Cancelamento:  make(chan struct{}),
	},
	{
		Nome:          "Crispy Turing",
		TempoCorte:    2 * time.Second,
		TempoGrelha:   10 * time.Second,
		TempoMontagem: 1 * time.Second,
		TempoBebida:   0,
		Cancelamento:  make(chan struct{}),
	},
	{
		Nome:          "Mongo Melt",
		TempoCorte:    1 * time.Second,
		TempoGrelha:   3 * time.Second,
		TempoMontagem: 0,
		TempoBebida:   0,
		Cancelamento:  make(chan struct{}),
	},
	{
		Nome:          "Float Juice",
		TempoCorte:    4 * time.Second,
		TempoGrelha:   0,
		TempoMontagem: 0,
		TempoBebida:   3 * time.Second,
		Cancelamento:  make(chan struct{}),
	},
	{
		Nome:          "Async Berry",
		TempoCorte:    2 * time.Second,
		TempoGrelha:   0,
		TempoMontagem: 0,
		TempoBebida:   2 * time.Second,
		Cancelamento:  make(chan struct{}),
	},
}

func StartWorkers(wg *sync.WaitGroup) {
	go workers.CortarWorker.Cortar(wg)
	go workers.GrelharWorker.Grelhar(wg)
	go workers.MontarWorker.Montar(wg)
	go workers.BebidaWorker.PrepararBebida(wg)
}

func DispatchPedidos(wg *sync.WaitGroup) {
	for _, pedido := range Menu {
		etapas := 0

		// Contabiliza as etapas que precisam ser realizadas
		if pedido.TempoCorte > 0 {
			etapas++
		}
		if pedido.TempoGrelha > 0 {
			etapas++
		}
		if pedido.TempoMontagem > 0 {
			etapas++
		}
		if pedido.TempoBebida > 0 {
			etapas++
		}

		pedido.QuantidadeTarefas = etapas
		wg.Add(pedido.QuantidadeTarefas)

		fmt.Printf("Novo pedido recebido: %s\n", pedido.Nome)

		if pedido.TempoCorte > 0 {
			workers.CortarWorker.Tarefas <- pedido
		} else if pedido.TempoGrelha > 0 {
			workers.GrelharWorker.Tarefas <- pedido
		} else if pedido.TempoMontagem > 0 {
			workers.MontarWorker.Tarefas <- pedido
		} else if pedido.TempoBebida > 0 {
			workers.BebidaWorker.Tarefas <- pedido
		}
	}
}

func CancelarPedido(nomePedido string) {
	for i := range Menu {
		if Menu[i].Nome == nomePedido {
			close(Menu[i].Cancelamento)  // Sinaliza o cancelamento
			fmt.Printf("%sPedido %s foi cancelado.%s\n", Vermelho, nomePedido, Branco)
			return
		}
	}
	fmt.Printf("%sPedido %s não encontrado.%s\n", Vermelho, nomePedido, Branco)
}

const (
	Branco   = "\033[0m"
	Vermelho = "\033[31m"
	Verde    = "\033[32m"
	Amarelo  = "\033[33m"
	Rosa     = "\033[35m"
	Ciana    = "\033[36m"
)

