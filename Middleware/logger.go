package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		fmt.Printf("Status:%d | Time:%v | Path:%s\n",
			c.Writer.Status(),
			duration,
			c.Request.URL.Path,
		)
	}
}
