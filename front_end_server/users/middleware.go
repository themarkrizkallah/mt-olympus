package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"front_end_server/common"
)


func AuthRequired() func(c *gin.Context) {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No authorization"})
			return
		}

		redisClient := common.GetRedisClient()

		_, err = redisClient.Get(cookieValue).Result()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization expired"})
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
	}
}

