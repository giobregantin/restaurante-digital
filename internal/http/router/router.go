package router

import (
	"context"

	_ "github.com/lib/pq"

	"github.com/hsxflowers/restaurante-digital/config"
	orderHandler "github.com/hsxflowers/restaurante-digital/internal/http/order"
	"github.com/hsxflowers/restaurante-digital/order"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Handlers(envs *config.Environments) *echo.Echo {
	e := echo.New()
	ctx := context.Background()

	// var orderDb domain.orderDatabase
	// var err error

	// db, err := sql.Open("postgres", "postgres://hsxflowers:N9tiQ81qNuzKP8Axv0Ae8aXZU8Pg6APf@dpg-cnl6bn7sc6pc73cc28vg-a.oregon-postgres.render.com/dborder?sslmode=require")
	// if err != nil {
	// 	log.Fatal("Error connecting to the database: ", err)
	// }

	// orderDb = orderDatabase.NewSQLStore(db)

	// log.Debug("")

	// orderRepository := order.NewOrderRepository(orderDb)
	orderService := order.NewOrderService()
	orderHandler := orderHandler.NewOrderHandler(ctx, orderService)

	e.GET("/swagger*", echoSwagger.WrapHandler)

	order := e.Group("order")
	order.POST("", orderHandler.Create)

	return e
}
