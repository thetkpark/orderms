package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pallat/micro/order"
)

// HandlerFunc defined new type for handler func
type HandlerFunc func(ctx order.Context)

type Router struct {
	*gin.Engine
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.Engine.POST(path, func(context *gin.Context) {
		handler(&Context{context})
	})
}

type Context struct {
	*gin.Context
}

func (c *Context) Order() (o order.Order, err error) {
	err = c.ShouldBindJSON(&o)
	return
}

func (c *Context) JSON(code int, v interface{}) {
	c.Context.JSON(code, v)
}
