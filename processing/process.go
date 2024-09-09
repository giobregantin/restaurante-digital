package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hsxflowers/restaurante-digital/processing/db"
	"github.com/hsxflowers/restaurante-digital/processing/domain"
	"github.com/hsxflowers/restaurante-digital/workers"
)

type Process struct {
	Wg             *sync.WaitGroup
	RestauranteeDb db.RestauranteDatabase
	ctx            context.Context
}

func NewProcess(wg *sync.WaitGroup, restauranteDb db.RestauranteDatabase, ctx context.Context) *Process {
	return &Process{
		Wg:             wg,
		RestauranteeDb: restauranteDb,
		ctx:            ctx,
	}
}

var Menu = []*domain.Pedido{
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "123",
		ItemId:       "Callback Burguer",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     "abcd",
		UsuarioId:    "1234",
		ItemId:       "Null-Burguer",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "456",
		ItemId:       "Crispy Turing",
		Cancelamento: make(chan struct{}),
		Prioridade: true,
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "456",
		ItemId:       "Mongo Melt",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "789",
		ItemId:       "Webwrap",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "789",
		ItemId:       "NPM Nuggets",
		Cancelamento: make(chan struct{}),
		Prioridade: true,
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "789",
		ItemId:       "Webwrap",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "1011",
		ItemId:       "Float Juice",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "789",
		ItemId:       "Array Apple",
		Cancelamento: make(chan struct{}),
	},
	{
		PedidoId:     uuid.New().String(),
		UsuarioId:    "789",
		ItemId:       "Async Berry",
		Cancelamento: make(chan struct{}),
	},
}

func (p *Process) StartWorkers() {
	go workers.CortarWorker.Cortar(p.ctx, p.Wg, p.RestauranteeDb)
	go workers.GrelharWorker.Grelhar(p.ctx, p.Wg, p.RestauranteeDb)
	go workers.MontarWorker.Montar(p.ctx, p.Wg, p.RestauranteeDb)
	go workers.BebidaWorker.PrepararBebida(p.ctx, p.Wg, p.RestauranteeDb)
}

func (p *Process) DispatchPedidos(ctx context.Context) {
	prioritarios := []*domain.Pedido{}
	naoPrioritarios := []*domain.Pedido{}

	for _, pedido := range Menu {
		if pedido.Prioridade {
			prioritarios = append(prioritarios, pedido)
		} else {
			naoPrioritarios = append(naoPrioritarios, pedido)
		}
	}

	ordemPedidos := append(prioritarios, naoPrioritarios...)

	for i, pedido := range ordemPedidos {

		item, err := p.RestauranteeDb.GetItem(ctx, pedido.ItemId)
		if err != nil {
			fmt.Printf("%sFalha ao processar o pedido com ItemId %s. Pedido ignorado.%s\n.", Vermelho, pedido.ItemId, Branco)
			continue
		}

		pedido.Nome = item.Nome
		pedido.TempoCorte = item.TempoCorte
		pedido.TempoGrelha = item.TempoGrelha
		pedido.TempoMontagem = item.TempoMontagem
		pedido.TempoBebida = item.TempoBebida
		pedido.Valor = item.Valor
		pedido.Status = "Em andamento"

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

		tempoEstimado, err := p.CalcularTempoEstimado(ctx, pedido)
		if err != nil {
			fmt.Printf("Erro ao calcular o tempo estimado: %v\n", err)
			continue
		}

		pedido.QuantidadeTarefas = etapas
		pedido.TempoEstimado = tempoEstimado
		Menu[i] = pedido

		p.Wg.Add(pedido.QuantidadeTarefas)
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

		err = p.RestauranteeDb.CreatePedido(ctx, pedido)
		if err != nil {
			fmt.Printf("%sErro ao adicionar o pedido com ItemId %s no banco de dados.%s\n.", Vermelho, pedido.ItemId, Branco)
			continue
		}

	}
}

func (p *Process) CalcularTempoEstimado(ctx context.Context, pedidoAtual *domain.Pedido) (time.Duration, error) {
	tempoEstimado := time.Duration(0)

	pedidosAnteriores, err := p.RestauranteeDb.GetPedidosAnteriores(ctx, pedidoAtual.PedidoId)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar pedidos anteriores: %w", err)
	}

	for _, pedidoAnterior := range pedidosAnteriores {
		if pedidoAnterior.Status == "Em andamento" {
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
	}

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

	return tempoEstimado, nil
}

func CancelarPedido(ctx context.Context, pedidoId string, db db.RestauranteDatabase) error {
	for i, pedido := range Menu {
		if pedido.PedidoId == pedidoId {
			close(Menu[i].Cancelamento)
			fmt.Printf("%sPedido com ID %s foi cancelado.%s\n", Vermelho, pedidoId, Branco)

			err := db.UpdatePedidoStatus(ctx, pedidoId, "Cancelado")
			if err != nil {
				fmt.Printf("%sErro ao atualizar o status no banco de dados para o pedido com ID %s.%s\n", Vermelho, pedidoId, Branco)
				return err
			}
			return nil
		}
	}

	fmt.Printf("%sPedido com ID %s não encontrado.%s\n", Vermelho, pedidoId, Branco)
	return fmt.Errorf("pedido com ID %s não encontrado", pedidoId)
}

func PedirConta(ctx context.Context, restauranteDb db.RestauranteDatabase, usuarioId string) {
    pedidos, valorTotal, err := restauranteDb.GetPedidos(ctx, usuarioId)
    if err != nil {
        fmt.Printf("Erro ao obter a conta: %s\n", err)
        return
    }

	fmt.Println("=====================================================")
	fmt.Println("Conta do usuário:", usuarioId)
    fmt.Println("Itens pedidos:")
    for _, pedido := range pedidos {
        fmt.Printf("- %s: R$ %.2f\n", pedido.Nome, pedido.Valor)
    }
    fmt.Printf("\nValor total: R$ %.2f\n", valorTotal)
	fmt.Println("=====================================================")
}

func DeletarPedidos(ctx context.Context, db db.RestauranteDatabase) error {

	err := db.DeletarPedidos(ctx)
	if err != nil {
		fmt.Printf("%sErro ao reiniciar tabela pedido.%s\n", Vermelho, Branco)
		return err
	}
	fmt.Printf("%sTabela pedido foi reiniciada.%s\n", Amarelo, Branco)
	return nil

}

const (
	Branco   = "\033[0m"
	Vermelho = "\033[31m"
	Verde    = "\033[32m"
	Amarelo  = "\033[33m"
	Rosa     = "\033[35m"
	Ciana    = "\033[36m"
)
