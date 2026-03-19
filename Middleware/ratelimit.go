package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func Ratelimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(1), 5)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "limit exceded"})
		}
	}
}
