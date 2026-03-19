package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Generatejwt(userId uint64) (string, error) {
	var jwtKey = []byte(os.Getenv("KEY"))
	fmt.Println(" KEY:", os.Getenv("KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userId,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)

}
