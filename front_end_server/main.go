package main

import (
	"github.com/gin-gonic/gin"

	"front_end_server/client"
	"front_end_server/env"
	"front_end_server/order"
)

func startServer(){
	client.InitExchangeService()
	defer client.Cleanup()

	r := gin.Default()

	r.POST("/orders/", order.CreateOrder)

	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	env.Init()
	startServer()
}
