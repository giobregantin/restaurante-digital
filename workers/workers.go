package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hsxflowers/restaurante-digital/processing/db"
	"github.com/hsxflowers/restaurante-digital/processing/domain"
)

type Worker struct {
	Nome    string
	Tarefas chan *domain.Pedido
	Context context.Context
}

var (
	CortarWorker  = Worker{Nome: "Cortar", Tarefas: make(chan *domain.Pedido, 20)}
	GrelharWorker = Worker{Nome: "Grelhar", Tarefas: make(chan *domain.Pedido, 20)}
	MontarWorker  = Worker{Nome: "Montar", Tarefas: make(chan *domain.Pedido, 20)}
	BebidaWorker  = Worker{Nome: "Bebida", Tarefas: make(chan *domain.Pedido, 20)}
)

func (w *Worker) Cortar(ctx context.Context, wg *sync.WaitGroup, db db.RestauranteDatabase) {
	for pedido := range w.Tarefas {
		select {
		case <-pedido.Cancelamento:
			fmt.Printf("[%s] %sPedido %s cancelado durante o corte.%s\n", w.Nome, Vermelho, pedido.Nome, Branco)
			for i := 0; i < pedido.QuantidadeTarefas; i++ {
				wg.Done()
			}
			continue
		default:
			fmt.Printf("[%s] %sIniciando corte para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
			time.Sleep(pedido.TempoCorte)
			fmt.Printf("[%s] %sConcluído o corte para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
			if pedido.TempoGrelha > 0 {
				GrelharWorker.Tarefas <- pedido
			} else if pedido.TempoMontagem > 0 {
				MontarWorker.Tarefas <- pedido
			} else if pedido.TempoBebida > 0 {
				BebidaWorker.Tarefas <- pedido
			} else {
				fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
				err := db.UpdatePedidoStatus(ctx, pedido.PedidoId, "Concluído")
				if err != nil {
					fmt.Printf("%sErro ao atualizar o status no banco de dados para o pedido com ID %s.%s\n", Vermelho, pedido.PedidoId, Branco)
				}
			}
			wg.Done()
		}
	}
}

func (w *Worker) Grelhar(ctx context.Context, wg *sync.WaitGroup, db db.RestauranteDatabase) {
	for pedido := range w.Tarefas {
		select {
		case <-pedido.Cancelamento:
			fmt.Printf("[%s] %sPedido %s cancelado durante a grelha.%s\n", w.Nome, Vermelho, pedido.Nome, Branco)
			for i := 0; i < pedido.QuantidadeTarefas-1; i++ {
				wg.Done()
			}
			continue
		default:
			fmt.Printf("[%s] %sIniciando grelha para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
			time.Sleep(pedido.TempoGrelha)
			fmt.Printf("[%s] %sConcluído a grelha para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
			if pedido.TempoMontagem > 0 {
				MontarWorker.Tarefas <- pedido
			} else if pedido.TempoBebida > 0 {
				BebidaWorker.Tarefas <- pedido
			} else {
				fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
				err := db.UpdatePedidoStatus(ctx, pedido.PedidoId, "Concluído")
				if err != nil {
					fmt.Printf("%sErro ao atualizar o status no banco de dados para o pedido com ID %s.%s\n", Vermelho, pedido.PedidoId, Branco)
				}
			}
			wg.Done()
		}
	}
}

func (w *Worker) Montar(ctx context.Context, wg *sync.WaitGroup, db db.RestauranteDatabase) {
	for pedido := range w.Tarefas {
		select {
		case <-pedido.Cancelamento:
			fmt.Printf("[%s] %sPedido %s cancelado durante a montagem.%s\n", w.Nome, Vermelho, pedido.Nome, Branco)
			for i := 0; i < pedido.QuantidadeTarefas-2; i++ {
				wg.Done()
			}
			continue
		default:
			fmt.Printf("[%s] %sIniciando montagem para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
			time.Sleep(pedido.TempoMontagem)
			fmt.Printf("[%s] %sConcluído a montagem para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
			if pedido.TempoBebida > 0 {
				BebidaWorker.Tarefas <- pedido
			} else {
				fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
				err := db.UpdatePedidoStatus(ctx, pedido.PedidoId, "Concluído")
				if err != nil {
					fmt.Printf("%sErro ao atualizar o status no banco de dados para o pedido com ID %s.%s\n", Vermelho, pedido.PedidoId, Branco)
				}
			}
			wg.Done()
		}
	}
}

func (w *Worker) PrepararBebida(ctx context.Context, wg *sync.WaitGroup, db db.RestauranteDatabase) {
	for pedido := range w.Tarefas {
		select {
		case <-pedido.Cancelamento:
			fmt.Printf("[%s] %sPedido %s cancelado durante a preparação da bebida.%s\n", w.Nome, Vermelho, pedido.Nome, Branco)
			wg.Done()
			continue
		default:
			fmt.Printf("[%s] %sIniciando preparação da bebida para: %s%s\n", w.Nome, Ciana, pedido.Nome, Branco)
			time.Sleep(pedido.TempoBebida)
			fmt.Printf("[%s] %sConcluído a preparação da bebida para: %s%s\n", w.Nome, Verde, pedido.Nome, Branco)
			fmt.Printf("%sPedido %s finalizado com sucesso!%s\n", Rosa, pedido.Nome, Branco)
			err := db.UpdatePedidoStatus(ctx, pedido.PedidoId, "Concluído")
			if err != nil {
				fmt.Printf("%sErro ao atualizar o status no banco de dados para o pedido com ID %s.%s\n", Vermelho, pedido.PedidoId, Branco)
			}
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
