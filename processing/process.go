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
		Status:        "Iniciado",
	},
	{
		Nome:          "Null-Burguer",
		TempoCorte:    4 * time.Second,
		TempoGrelha:   7 * time.Second,
		TempoMontagem: 2 * time.Second,
		TempoBebida:   0,
		Status:        "Iniciado",
	},
	{
		Nome:          "Crispy Turing",
		TempoCorte:    2 * time.Second,
		TempoGrelha:   10 * time.Second,
		TempoMontagem: 1 * time.Second,
		TempoBebida:   0,
		Status:        "Iniciado",
	},
	{
		Nome:          "Mongo Melt",
		TempoCorte:    1 * time.Second,
		TempoGrelha:   3 * time.Second,
		TempoMontagem: 0,
		TempoBebida:   0,
		Status:        "Iniciado",
	},
	{
		Nome:          "Float Juice",
		TempoCorte:    4 * time.Second,
		TempoGrelha:   0,
		TempoMontagem: 0,
		TempoBebida:   3 * time.Second,
		Status:        "Iniciado",
	},
	{
		Nome:          "Async Berry",
		TempoCorte:    2 * time.Second,
		TempoGrelha:   0,
		TempoMontagem: 0,
		TempoBebida:   2 * time.Second,
		Status:        "Iniciado",
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

		wg.Add(etapas)

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

