package workers

import (
	"fmt"
	"sync"
	"time"

	"github.com/hsxflowers/restaurante-digital/processing/db"
	"github.com/hsxflowers/restaurante-digital/workers"
)

type Process struct {
	Wg             *sync.WaitGroup
	RestauranteeDb db.CatDatabase
}

func NewProcess(wg *sync.WaitGroup, restauranteeDb db.CatDatabase) domain.Service {
	return &Process{
		Wg:             wg,
		RestauranteeDb: restauranteeDb,
	}
}

var Menu = []workers.Pedido{
	{
		UsuarioId:    "123",
		ItemId:       "a1b2c3d4",
		Cancelamento: make(chan struct{}),
	},
	{
		UsuarioId:    "123",
		ItemId:       "e5f6g7h8",
		Cancelamento: make(chan struct{}),
	},
	{
		UsuarioId:    "456",
		ItemId:       "i9j0k1l2",
		Cancelamento: make(chan struct{}),
	},
	{
		UsuarioId:    "456",
		ItemId:       "m3n4o5p6",
		Cancelamento: make(chan struct{}),
	},
	{
		UsuarioId:    "789",
		ItemId:       "q7r8s9t0",
		Cancelamento: make(chan struct{}),
	},
	{
		UsuarioId:    "789",
		ItemId:       "u1v2w3x4",
		Cancelamento: make(chan struct{}),
	},
}

func (p *Process) StartWorkers(wg *sync.WaitGroup) {
	go workers.CortarWorker.Cortar(wg)
	go workers.GrelharWorker.Grelhar(wg)
	go workers.MontarWorker.Montar(wg)
	go workers.BebidaWorker.PrepararBebida(wg)
}

func (p *Process) DispatchPedidos(wg *sync.WaitGroup) {
	for i, pedido := range Menu {

		etapas := 0

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

		tempoEstimado := CalcularTempoEstimado(i)

		pedido.QuantidadeTarefas = etapas
		pedido.TempoEstimado = tempoEstimado
		Menu[i] = pedido

		wg.Add(pedido.QuantidadeTarefas)
		fmt.Printf("Novo pedido recebido: %s (Tempo estimado: %v)\n", pedido.Nome, pedido.TempoEstimado)

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

func CalcularTempoEstimado(index int) time.Duration {
	tempoEstimado := time.Duration(0)

	for j := 0; j < index; j++ {
		pedidoAnterior := Menu[j]
		pedidoAtual := Menu[index]

		if pedidoAnterior.TempoCorte > 0 && pedidoAtual.TempoCorte > 0 {
			tempoEstimado += pedidoAnterior.TempoCorte
		}
		if pedidoAnterior.TempoGrelha > 0 && pedidoAtual.TempoGrelha > 0 {
			tempoEstimado += pedidoAnterior.TempoGrelha
		}
		if pedidoAnterior.TempoMontagem > 0 && pedidoAtual.TempoMontagem > 0 {
			tempoEstimado += pedidoAnterior.TempoMontagem
		}
		if pedidoAnterior.TempoBebida > 0 && pedidoAtual.TempoBebida > 0 {
			tempoEstimado += pedidoAnterior.TempoBebida
		}
	}

	pedidoAtual := Menu[index]
	if pedidoAtual.TempoCorte > 0 {
		tempoEstimado += pedidoAtual.TempoCorte
	}
	if pedidoAtual.TempoGrelha > 0 {
		tempoEstimado += pedidoAtual.TempoGrelha
	}
	if pedidoAtual.TempoMontagem > 0 {
		tempoEstimado += pedidoAtual.TempoMontagem
	}
	if pedidoAtual.TempoBebida > 0 {
		tempoEstimado += pedidoAtual.TempoBebida
	}

	return tempoEstimado
}

func CancelarPedido(nomePedido string) {
	for i := range Menu {
		if Menu[i].Nome == nomePedido {
			close(Menu[i].Cancelamento)
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
