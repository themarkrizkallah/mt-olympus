package product

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/database"
)

func GetProducts(c *gin.Context) {
	products, err := database.GetProducts()
	if err != nil {
		log.Println("Error retrieving products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	c.JSON(http.StatusOK, products)
}
