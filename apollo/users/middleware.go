package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/redis"
)

const cookieName = "exchange_userCookie"

func AuthRequired() func(c *gin.Context) {
	return func(c *gin.Context) {
		var userId string

		sessionId, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No authorization"})
			return
		}

		if err := redis.Codec.GetContext(c, sessionId, &userId); err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization missing or expired"})
		} else if err != nil {
			log.Println("Error getting value from redis", err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "An error occurred"})
		} else {
			c.Set("user_id", userId)
		}
	}
}
