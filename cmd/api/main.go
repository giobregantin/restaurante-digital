package main

import (
	"fmt"
	"sync"

	process "github.com/hsxflowers/restaurante-digital/processing"
)

func main() {
	var wg sync.WaitGroup

	process.StartWorkers(&wg)
	process.DispatchOrders(&wg)

	// time.Sleep(10 * time.Second)
	// process.CancelarOrder("Crispy Turing")

	wg.Wait()
	fmt.Println("Todos os orders foram processados ou cancelados.")
}
