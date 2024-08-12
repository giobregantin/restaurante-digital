package main

import (
	"fmt"
	"sync"
	process "github.com/hsxflowers/restaurante-digital/processing"
)

func main() {
	var wg sync.WaitGroup

	process.StartWorkers(&wg)
	process.DispatchPedidos(&wg)

	wg.Wait()
	fmt.Println("Todos os pedidos foram processados.")
}
