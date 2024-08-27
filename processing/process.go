package workers

import (
	"fmt"
	"sync"
	"time"

	"github.com/hsxflowers/restaurante-digital/internal/workers"
)

var Menu = []workers.Order{
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

func DispatchOrders(wg *sync.WaitGroup) {
	for i, order := range Menu {
		etapas := 0

		if order.TempoCorte > 0 {
			etapas++
		}
		if order.TempoGrelha > 0 {
			etapas++
		}
		if order.TempoMontagem > 0 {
			etapas++
		}
		if order.TempoBebida > 0 {
			etapas++
		}

		tempoEstimado := CalcularTempoEstimado(i)

		order.QuantidadeTarefas = etapas
		order.TempoEstimado = tempoEstimado
		Menu[i] = order

		wg.Add(order.QuantidadeTarefas)
		fmt.Printf("Novo order recebido: %s (Tempo estimado: %v)\n", order.Nome, order.TempoEstimado)

		if order.TempoCorte > 0 {
			workers.CortarWorker.Tarefas <- order
		} else if order.TempoGrelha > 0 {
			workers.GrelharWorker.Tarefas <- order
		} else if order.TempoMontagem > 0 {
			workers.MontarWorker.Tarefas <- order
		} else if order.TempoBebida > 0 {
			workers.BebidaWorker.Tarefas <- order
		}
	}
}

func CalcularTempoEstimado(index int) time.Duration {
	tempoEstimado := time.Duration(0)

	for j := 0; j < index; j++ {
		orderAnterior := Menu[j]
		orderAtual := Menu[index]

		if orderAnterior.TempoCorte > 0 && orderAtual.TempoCorte > 0 {
			tempoEstimado += orderAnterior.TempoCorte
		}
		if orderAnterior.TempoGrelha > 0 && orderAtual.TempoGrelha > 0 {
			tempoEstimado += orderAnterior.TempoGrelha
		}
		if orderAnterior.TempoMontagem > 0 && orderAtual.TempoMontagem > 0 {
			tempoEstimado += orderAnterior.TempoMontagem
		}
		if orderAnterior.TempoBebida > 0 && orderAtual.TempoBebida > 0 {
			tempoEstimado += orderAnterior.TempoBebida
		}
	}

	orderAtual := Menu[index]
	if orderAtual.TempoCorte > 0 {
		tempoEstimado += orderAtual.TempoCorte
	}
	if orderAtual.TempoGrelha > 0 {
		tempoEstimado += orderAtual.TempoGrelha
	}
	if orderAtual.TempoMontagem > 0 {
		tempoEstimado += orderAtual.TempoMontagem
	}
	if orderAtual.TempoBebida > 0 {
		tempoEstimado += orderAtual.TempoBebida
	}

	return tempoEstimado
}

func CancelarOrder(nomeOrder string) {
	for i := range Menu {
		if Menu[i].Nome == nomeOrder {
			close(Menu[i].Cancelamento)
			fmt.Printf("%sOrder %s foi cancelado.%s\n", Vermelho, nomeOrder, Branco)
			return
		}
	}
	fmt.Printf("%sOrder %s não encontrado.%s\n", Vermelho, nomeOrder, Branco)
}

const (
	Branco   = "\033[0m"
	Vermelho = "\033[31m"
	Verde    = "\033[32m"
	Amarelo  = "\033[33m"
	Rosa     = "\033[35m"
	Ciana    = "\033[36m"
)
