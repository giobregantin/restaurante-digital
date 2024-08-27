package order

import (
	"context"
	"net/http"

	_ "github.com/hsxflowers/restaurante-digital/order"
	"github.com/hsxflowers/restaurante-digital/order/domain"

	"github.com/hsxflowers/restaurante-digital/exceptions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type OrderHandler struct {
	ctx          context.Context
	orderService domain.Service
}

func NewOrderHandler(ctx context.Context, orderService domain.Service) OrderHandler {
	return OrderHandler{
		ctx,
		orderService,
	}
}

// Get
//
//	@Summary		Ver item
//	@Description	Endpoint que permite a chamada de um item de acordo com o id informado.
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	[]domain.Item	"OK"
//	@Failure		404			"Item with id {:id} not found"
//	@Failure		500			"internal Server Error"
//	@Router			/order/{:id} [get]

// func (h *OrderHandler) GetItem(c echo.Context) error {
// 	itemId := c.Param("item_id")
// 	if itemId == "" || itemId == ":item_id" {
// 		log.Error("handler_get: item id is required.", exceptions.ErrIdIsNotValid)
// 		return exceptions.New(exceptions.ErrTagIsRequired, nil)
// 	}

// 	response, err := h.orderService.GetItem(h.ctx, itemId)
// 	if err != nil {
// 		log.Error("handler_get: error on get a item.", err)
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, response)
// }

// Create
//
//	@Summary		Criação de orders.
//	@Description	Endpoint que permite a criação de gatos.
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.OrderRequest	true	"body"
//	@Success		201		{object}	domain.OrderResponse
//	@Failure		422		"Unprocessable Json: Payload enviado com erro de syntax do json"
//	@Failure		400		"Erros de validação ou request indevido"
//	@Failure		500			"internal Server Error"
//	@Router			/channels [post]
func (h *OrderHandler) Create(c echo.Context) error {
	req := new(domain.OrderRequest)

	if err := c.Bind(req); err != nil {
		log.Error("handler_create: error marshal order", err)
		return exceptions.New(exceptions.ErrBadData, err)
	}

	if err := req.Validate(); err != nil {
		log.Error("handler_create: error on create order", err)
		return exceptions.New(err, nil)
	}

	order, err := h.orderService.CreateOrder(h.ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, order)
}
