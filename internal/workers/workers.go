package workers

import (
	"fmt"
	"sync"
	"time"
)

type Order struct {
	Nome              string
	TempoCorte        time.Duration
	TempoGrelha       time.Duration
	TempoMontagem     time.Duration
	TempoBebida       time.Duration
	QuantidadeTarefas int
	Cancelamento      chan struct{}
	TempoEstimado     time.Duration
}

type Worker struct {
	Nome    string
	Tarefas chan Order
}

var (
	CortarWorker  = Worker{Nome: "Cortar", Tarefas: make(chan Order, 20)}
	GrelharWorker = Worker{Nome: "Grelhar", Tarefas: make(chan Order, 20)}
	MontarWorker  = Worker{Nome: "Montar", Tarefas: make(chan Order, 20)}
	BebidaWorker  = Worker{Nome: "Bebida", Tarefas: make(chan Order, 20)}
)

func (w *Worker) Cortar(wg *sync.WaitGroup) {
	for order := range w.Tarefas {
		select {
		case <-order.Cancelamento:
			fmt.Printf("[%s] %sOrder %s cancelado durante o corte.%s\n", w.Nome, Vermelho, order.Nome, Branco)
			for i := 0; i < order.QuantidadeTarefas; i++ {
				wg.Done()
			}
			continue
		default:
			fmt.Printf("[%s] %sIniciando corte para: %s%s\n", w.Nome, Ciana, order.Nome, Branco)
			time.Sleep(order.TempoCorte)
			fmt.Printf("[%s] %sConcluído o corte para: %s%s\n", w.Nome, Verde, order.Nome, Branco)
			if order.TempoGrelha > 0 {
				GrelharWorker.Tarefas <- order
			} else if order.TempoMontagem > 0 {
				MontarWorker.Tarefas <- order
			} else if order.TempoBebida > 0 {
				BebidaWorker.Tarefas <- order
			} else {
				fmt.Printf("%sOrder %s finalizado com sucesso!%s\n", Rosa, order.Nome, Branco)
			}
			wg.Done()
		}
	}
}

func (w *Worker) Grelhar(wg *sync.WaitGroup) {
	for order := range w.Tarefas {
		select {
		case <-order.Cancelamento:
			fmt.Printf("[%s] %sOrder %s cancelado durante a grelha.%s\n", w.Nome, Vermelho, order.Nome, Branco)
			for i := 0; i < order.QuantidadeTarefas-1; i++ {
				wg.Done()
			}
			continue
		default:
			fmt.Printf("[%s] %sIniciando grelha para: %s%s\n", w.Nome, Ciana, order.Nome, Branco)
			time.Sleep(order.TempoGrelha)
			fmt.Printf("[%s] %sConcluído a grelha para: %s%s\n", w.Nome, Verde, order.Nome, Branco)
			if order.TempoMontagem > 0 {
				MontarWorker.Tarefas <- order
			} else if order.TempoBebida > 0 {
				BebidaWorker.Tarefas <- order
			} else {
				fmt.Printf("%sOrder %s finalizado com sucesso!%s\n", Rosa, order.Nome, Branco)
			}
			wg.Done()
		}
	}
}

func (w *Worker) Montar(wg *sync.WaitGroup) {
	for order := range w.Tarefas {
		select {
		case <-order.Cancelamento:
			fmt.Printf("[%s] %sOrder %s cancelado durante a montagem.%s\n", w.Nome, Vermelho, order.Nome, Branco)
			for i := 0; i < order.QuantidadeTarefas-2; i++ {
				wg.Done()
			}
			continue
		default:
			fmt.Printf("[%s] %sIniciando montagem para: %s%s\n", w.Nome, Ciana, order.Nome, Branco)
			time.Sleep(order.TempoMontagem)
			fmt.Printf("[%s] %sConcluído a montagem para: %s%s\n", w.Nome, Verde, order.Nome, Branco)
			if order.TempoBebida > 0 {
				BebidaWorker.Tarefas <- order
			} else {
				fmt.Printf("%sOrder %s finalizado com sucesso!%s\n", Rosa, order.Nome, Branco)
			}
			wg.Done()
		}
	}
}

func (w *Worker) PrepararBebida(wg *sync.WaitGroup) {
	for order := range w.Tarefas {
		select {
		case <-order.Cancelamento:
			fmt.Printf("[%s] %sOrder %s cancelado durante a preparação da bebida.%s\n", w.Nome, Vermelho, order.Nome, Branco)
			wg.Done()
			continue
		default:
			fmt.Printf("[%s] %sIniciando preparação da bebida para: %s%s\n", w.Nome, Ciana, order.Nome, Branco)
			time.Sleep(order.TempoBebida)
			fmt.Printf("[%s] %sConcluído a preparação da bebida para: %s%s\n", w.Nome, Verde, order.Nome, Branco)
			fmt.Printf("%sOrder %s finalizado com sucesso!%s\n", Rosa, order.Nome, Branco)
			wg.Done()
		}
	}
}

const (
	Branco   = "\033[0m"
	Vermelho = "\033[31m"
	Verde    = "\033[32m"
	Amarelo  = "\033[33m"
	Rosa     = "\033[35m"
	Ciana    = "\033[36m"
)
