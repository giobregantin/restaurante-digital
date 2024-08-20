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
		Cancelamento:  make(chan struct{}),
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

