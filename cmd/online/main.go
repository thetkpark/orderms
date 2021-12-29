package main

import (
	"github.com/joho/godotenv"
	"github.com/pallat/micro/order"
	"github.com/pallat/micro/router"
	store2 "github.com/pallat/micro/store"
	"log"
	"os"
)

func init() {
	err := godotenv.Load("online.env")
	if err != nil {
		log.Fatalf("please consider environment variable: %s\n", err)
	}
}

func main() {
	r := router.NewRouter()

	s := store2.NewMongoDBStore(os.Getenv("DSN"))
	handler := order.NewHandler(os.Getenv("FILTER_CHANNEL"), s)
	r.POST("/api/v1/orders", handler.Order)

	r.ListenAndServe()()
}
