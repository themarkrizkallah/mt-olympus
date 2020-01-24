package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zeus/redis"
)

func AuthRequired() func(c *gin.Context) {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No authorization"})
			return
		}

		_, err = redis.GetClient().Get(cookieValue).Result()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization expired"})
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
	}
}
