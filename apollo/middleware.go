package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v7"

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

		err = redis.Codec.GetContext(c, sessionId, &userId)
		if err != nil {
			if errors.Is(err, cache.ErrCacheMiss) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization missing or expired"})
			} else {
				log.Println("Error getting value from redis:", err)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "An error occurred"})
			}

			return
		}

		c.Set("user_id", userId)
	}
}
