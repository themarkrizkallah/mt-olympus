package order

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"apollo/kafka"
	pb "apollo/proto"
	"apollo/redis"
)

const userIdKey = "user_id"

func PostOrder(c *gin.Context) {
	var (
		err error
		request pb.OrderRequest
	)

	userId := c.GetString(userIdKey)

	if err := c.BindJSON(&request); err != nil {
		log.Println("Error parsing payload:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "bad request"})
		return
	}

	// Validate the order request
	if request.GetAmount() <= 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "amount must be > 0"})
		return
	}
	if request.GetPrice() <= 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "price must be > 0"})
		return
	}
	if _, ok := pb.Side_value[request.GetSide().String()]; !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "side must be 0 or 1"})
		return
	}
	if _, ok := pb.Type_value[request.GetType().String()]; !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "type must be 0, 1, or 2"})
		return
	}

	if productsMap, err := redis.GetProductsMap(c); err != nil {
		log.Println("Error retrieving products", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	} else if _, ok := productsMap[request.GetProductId()]; !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid product_id"})
		return
	}

	request.UserId = userId
	request.OrderId = uuid.New().String()
	log.Println("Order Router - processing order request\n")

	confChan, err := kafka.SendOrderRequest(request)
	if err != nil {
		log.Fatalln("Order Router - Error sending request to kafka:", err)
	}
	c.JSON(http.StatusOK, newConf(<-confChan))
}
