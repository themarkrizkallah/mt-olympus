package order

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	"apollo/kafka"
	pb "apollo/proto"
	"apollo/redis"
)

const (
	userIdKey         = "user_id"
	orderRequestTopic = "order.request"
)

func PostOrder(c *gin.Context) {
	var orderRequest pb.OrderRequest

	userId := c.GetString(userIdKey)

	if err := c.BindJSON(&orderRequest); err != nil {
		log.Println("Error parsing payload:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not parse payload"})
		return
	}

	// Validate the order request
	if orderRequest.GetAmount() <= 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "amount must be > 0"})
		return
	} else if orderRequest.GetPrice() <= 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "price must be > 0"})
		return
	} else if _, ok := pb.Side_value[orderRequest.GetSide().String()]; !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "side must be 0 or 1"})
		return
	} else if _, ok := pb.Type_value[orderRequest.GetType().String()]; !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "type must be 0, 1, or 2"})
		return
	}

	productsMap, err := redis.GetProductsMap(c)
	if err != nil {
		log.Println("Error retrieving products", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	} else if _, ok := productsMap[orderRequest.GetProductId()]; !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid product_id"})
		return
	}

	orderRequest.UserId = userId
	orderRequest.OrderId = uuid.New().String()
	log.Printf("Processing: %+v\n", orderRequest)

	data, err := proto.Marshal(&orderRequest)
	if err != nil {
		log.Println("Marshaling error: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	sendOp := kafka.SendOp{
		ID:       orderRequest.OrderId,
		Topic:    orderRequestTopic,
		Value:    data,
		Receiver: make(chan pb.OrderConf),
	}

	kafka.Sender <- sendOp

	// Get OrderConf and only keep relevant fields
	orderConf := <-sendOp.Receiver
	orderConf.UserId = ""

	c.JSON(http.StatusOK, orderConf)
}
