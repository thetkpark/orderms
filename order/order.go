package order

import (
	"fmt"
	"net/http"
)

type Context interface {
	Order() (Order, error)
	JSON(code int, v interface{})
}

type Handler struct {
	channel string
}

func (h *Handler) Order(c Context) {
	order, err := c.Order()
	if err != nil {

	}

	if order.SalesChannel != h.channel {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("%s is not accepted", order.SalesChannel),
		})
	}
}
