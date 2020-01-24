package order

import (
	pb "apollo/proto"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	"apollo/kafka"
	"apollo/redis"
	"apollo/types"
)

const cookieName = "exchange_userCookie"

func CreateOrder(c *gin.Context) {
	var payload types.Payload

	sessionId, _ := c.Cookie(cookieName)
	redisClient := redis.GetClient()

	userId, err := redisClient.Get(sessionId).Result()
	if err == redis.Nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization expired or not logged in"})
	}

	if err = c.BindJSON(&payload); err != nil {
		log.Println("Error parsing payload:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not parse payload"})
		return
	}

	order := types.Order{
		UserId:    userId,
		OrderId:   uuid.New().String(),
		Amount:    payload.Amount,
		Price:     payload.Price,
		Side:      payload.Side,
		Type:      payload.Type,
		CreatedAt: time.Now(),
	}

	orderRequest := order.ToOrderRequest()
	log.Printf("Processing: %+v\n", orderRequest)

	data, err := proto.Marshal(&orderRequest)
	if err != nil {
		log.Println("Marshaling error: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	sendOp := kafka.SendOp{
		ID:       orderRequest.UserId,
		Topic:    "order.request",
		Value:    data,
		Receiver: make(chan pb.OrderConf),
	}

	//log.Println("sendOp:", sendOp)
	kafka.Sender <- sendOp
	c.JSON(http.StatusOK, gin.H{"response": types.FromProto(<-sendOp.Receiver)})
}
