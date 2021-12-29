package order

import (
	"fmt"
	"net/http"
)

type Context interface {
	Order() (Order, error)
	JSON(code int, v interface{})
}

type Store interface {
	Save(order Order) error
}

type Handler struct {
	channel string
	store   Store
}

func (h *Handler) Order(c Context) {
	order, err := c.Order()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if order.SalesChannel != h.channel {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("%s is not accepted", order.SalesChannel),
		})
		return
	}

	err = h.store.Save(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
}
