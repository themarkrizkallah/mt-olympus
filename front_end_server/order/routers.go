package order

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"

	"front_end_server/client"
	"front_end_server/common"
	"front_end_server/users"
)

func CreateOrder(c *gin.Context) {
	const cookieName = "exchange_userCookie"

	var (
		payload Payload
		user    users.User
	)

	cookieValue, _ := c.Cookie(cookieName)
	redisClient := common.GetRedisClient()

	id, err := redisClient.Get(cookieValue).Result()
	if err == redis.Nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization expired"})
	}

	err = c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err = users.FindUserById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order := payload.Parse()
	order.UserId = user.Id.Hex()
	order.OrderId = uuid.New().String()
	log.Printf("%+v\n", order)

	orderRequest := order.ToOrderRequest()
	exchangeServiceClient := client.GetExchangeServiceClient()

	resp, err := (*exchangeServiceClient).PlaceOrder(c, &orderRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"response": resp.Body})
	}
}
