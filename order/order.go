package order

type Context interface {
	Order() (Order, error)
	JSON(code int, v interface{})
}

type Handler struct {
	channel string
}

func (h *Handler) Order(c Context) {

}
