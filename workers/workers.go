// workers/workers.go
package workers

import (
	"fmt"
	"sync"
	"time"
)

type Pedido struct {
	Nome          string
	TempoCorte    time.Duration
	TempoGrelha   time.Duration
	TempoMontagem time.Duration
	TempoBebida   time.Duration
	Status        string
}

type Worker struct {
	Nome    string
	Tarefas chan Pedido
}

var (
	CortarWorker  = Worker{Nome: "Cortar", Tarefas: make(chan Pedido, 20)}
	GrelharWorker = Worker{Nome: "Grelhar", Tarefas: make(chan Pedido, 20)}
	MontarWorker  = Worker{Nome: "Montar", Tarefas: make(chan Pedido, 20)}
	BebidaWorker  = Worker{Nome: "Bebida", Tarefas: make(chan Pedido, 20)}
)

func (w *Worker) Cortar(wg *sync.WaitGroup) {
	for pedido := range w.Tarefas {
		fmt.Printf("[%s] %sIniciando corte para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
		time.Sleep(pedido.TempoCorte)
		pedido.Status = "Corte Concluído"
		fmt.Printf("[%s] %sConcluído o corte para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
		if pedido.TempoGrelha > 0 {
			GrelharWorker.Tarefas <- pedido
		} else if pedido.TempoMontagem > 0 {
			MontarWorker.Tarefas <- pedido
		} else if pedido.TempoBebida > 0 {
			BebidaWorker.Tarefas <- pedido
		} else {
			fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
		}
		wg.Done()
	}
}

func (w *Worker) Grelhar(wg *sync.WaitGroup) {
	for pedido := range w.Tarefas {
		fmt.Printf("[%s] %sIniciando grelha para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
		time.Sleep(pedido.TempoGrelha)
		pedido.Status = "Grelha Concluída"
		fmt.Printf("[%s] %sConcluído a grelha para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
		if pedido.TempoMontagem > 0 {
			MontarWorker.Tarefas <- pedido
		} else if pedido.TempoBebida > 0 {
			BebidaWorker.Tarefas <- pedido
		} else {
			fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
		}
		wg.Done()
	}
}

func (w *Worker) Montar(wg *sync.WaitGroup) {
	for pedido := range w.Tarefas {
		fmt.Printf("[%s] %sIniciando montagem para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
		time.Sleep(pedido.TempoMontagem)
		pedido.Status = "Montagem Concluída"
		fmt.Printf("[%s] %sConcluído a montagem para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
		if pedido.TempoBebida > 0 {
			BebidaWorker.Tarefas <- pedido
		} else {
			fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
		}
		wg.Done()
	}
}

func (w *Worker) PrepararBebida(wg *sync.WaitGroup) {
	for pedido := range w.Tarefas {
		fmt.Printf("[%s] %sIniciando preparação da bebida para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
		time.Sleep(pedido.TempoBebida)
		pedido.Status = "Bebida Concluída"
		fmt.Printf("[%s] %sConcluído a bebida para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
		fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
		wg.Done()
	}
}

const (
	Branco   = "\033[0m"
	Vermelho     = "\033[31m"
	Verde   = "\033[32m"
	Amarelo  = "\033[33m"
	Rosa = "\033[35m"
	Ciana    = "\033[36m"
)
