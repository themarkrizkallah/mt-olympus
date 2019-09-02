package order

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"front_end_server/client"
)

func CreateOrder(c *gin.Context) {
	var payload Payload
	err := c.BindJSON(&payload)

	if err != nil {
		panic(err)
	}

	order := payload.Parse()
	orderObj := order.ToProtoObj()
	exchangeServiceClient := client.GetExchangeServiceClient()

	if resp, err := (*exchangeServiceClient).CreateOrder(c, &orderObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"response": fmt.Sprint(resp.Response),
		})
	}
}
