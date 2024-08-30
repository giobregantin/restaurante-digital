package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	process "github.com/hsxflowers/restaurante-digital/processing"
	"github.com/hsxflowers/restaurante-digital/processing/db"
)

func main() {
	var wg sync.WaitGroup

	var restauranteDb db.CatDatabase
	var err error

	db, err := sql.Open("postgres", "postgresql://admin:Ep7elyAZpA5dYxZMe3vOjIIJJ1dXF3XZ@dpg-cr7fi923esus7388q3c0-a.oregon-postgres.render.com/dbrestaurante")
	if err != nil {
		log.Fatal("Erro ao conectar com DB: ", err)
	}

	restauranteDb = restauranteDatabase.NewSQLStore(db)

	process.StartWorkers(&wg)
	process.DispatchPedidos(&wg)

	time.Sleep(5 * time.Second)
	process.CancelarPedido("Crispy Turing")

	wg.Wait()
	fmt.Println("Todos os pedidos foram processados ou cancelados.")
}
