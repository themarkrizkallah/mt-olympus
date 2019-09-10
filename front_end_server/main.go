package main

import (
	"github.com/gin-gonic/gin"

	"front_end_server/client"
	"front_end_server/common"
	"front_end_server/env"
	"front_end_server/order"
	"front_end_server/users"
)

func startServer() {
	client.InitExchangeService()
	defer client.Cleanup()

	r := gin.Default()

	r.POST("/signup/", users.SignUp)
	r.POST("/login/", users.Login)

	// Endpoints that require authentication
	r.GET("/users/", users.AuthRequired(), users.ListUsers)
	r.POST("/orders/", users.AuthRequired(), order.CreateOrder)

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
