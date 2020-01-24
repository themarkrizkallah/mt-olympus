package order

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"zeus/client"
	"zeus/redis"
	"zeus/users"
)

func CreateOrder(c *gin.Context) {
	const cookieName = "exchange_userCookie"

	var (
		payload Payload
		user    users.User
	)

	cookieValue, _ := c.Cookie(cookieName)
	redisClient := redis.GetClient()

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
