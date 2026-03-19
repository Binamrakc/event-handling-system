package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Passwordhash() gin.HandlerFunc {
	return func(c *gin.Context) {
		Password := bcrypt.GenerateFromPassword([]byte(password))
	}
}
