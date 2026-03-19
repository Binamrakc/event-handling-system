package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unable to authorize"})
			c.Abort()
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		tokens, err := jwt.Parse(tokenString, func(tokens *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("KEY")), nil
		})
		if err != nil || !tokens.Valid {
			c.JSON(401, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		fmt.Println("GEN KEY:", os.Getenv("KEY"))
		c.Next()
	}
}
