package order

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"front_end_server/client"
)

func CreateOrder(c *gin.Context) {
	var payload Payload

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order := payload.Parse()
	orderRequest := order.ToOrderRequest()
	exchangeServiceClient := client.GetExchangeServiceClient()

	resp, err := (*exchangeServiceClient).PlaceOrder(c, &orderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"response": resp.Body})
	}
}
