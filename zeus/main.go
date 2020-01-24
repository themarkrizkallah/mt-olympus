package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"zeus/client"
	"zeus/database"
	"zeus/env"
	"zeus/order"
	"zeus/redis"
	"zeus/users"
)

/*
message OrderRequest {
    string user_id = 1;
    string order_id = 2;
    uint64 amount = 3;
    uint64 price = 4;
    Side Side = 5;
    Type type = 6;
    google.protobuf.Timestamp created_at = 7;
}
*/

const schema = `
{
    "type": "record",
    "name": "order",
    "fields": [
      { "name": "user_id", "type": "string" },
      { "name": "order_id", "type": "string" },
      { "name": "amount", "type": "long" },
      { "name": "price", "type": "long" },
      { "name": "side", "type": "boolean" },
      { "name": "created_at", "type": "string" }
    ]
}
`

func main() {
	env.Init()
	setupMongo()
	setupRedis()
	startServer()
}

func setupMongo() {
	_, err := database.Init()
	if err != nil {
		log.Fatalln("Error setting up Mongo:", err)
	}

	errs := database.SetupIndices(database.DefaultIndexConfig(), env.MongoDb)
	if len(errs) > 0 {
		log.Println("Errors setting up indices:", errs)
	}
}

func setupRedis() {
	_, err := redis.Init()
	if err != nil {
		log.Fatalln("Error setting up Redis:", err)
	}
}

func startServer() {
	client.InitExchangeService()
	defer client.Cleanup()

	r := gin.Default()

	r.POST("/signup/", users.SignUp)
	r.POST("/login/", users.Login)

	// Endpoints that require authentication
	authorized := r.Group("/")
	authorized.Use(users.AuthRequired())
	{
		authorized.GET("/users/", users.ListUsers)
		authorized.POST("/orders/", order.CreateOrder)
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}
