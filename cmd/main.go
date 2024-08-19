package main

import (
	"fmt"
	"sync"
	"time"

	process "github.com/hsxflowers/restaurante-digital/processing"
)

func main() {
	var wg sync.WaitGroup

	process.StartWorkers(&wg)
	process.DispatchPedidos(&wg)

	time.Sleep(10 * time.Second)
	process.CancelarPedido("Crispy Turing")

	wg.Wait()
	fmt.Println("Todos os pedidos foram processados ou cancelados.")
}
