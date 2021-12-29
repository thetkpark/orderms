package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pallat/micro/order"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

// HandlerFunc defined new type for handler func
type HandlerFunc func(ctx order.Context)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	r := gin.Default()
	return &Router{r}
}

func (r *Router) GET(path string, handler HandlerFunc) {
	r.Engine.GET(path, func(ctx *gin.Context) {
		handler(&Context{ctx})
	})
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.Engine.POST(path, func(ctx *gin.Context) {
		handler(&Context{ctx})
	})
}

func (r *Router) ListenAndServe() func() {
	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		stop()
		fmt.Println("shutting down gracefully, press Ctrl+C again to force")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(timeoutCtx); err != nil {
			fmt.Println(err)
		}
	}
}
