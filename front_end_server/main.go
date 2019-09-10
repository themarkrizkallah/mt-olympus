package main

import (
	"front_end_server/common"
	"front_end_server/users"
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
	r.POST("/signup/", users.SignUp)
	r.POST("/login/", users.Login)
	r.GET("/users/", users.ListUsers)

	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	env.Init()
	common.InitMongo()
	common.InitRedis()
	startServer()
}
