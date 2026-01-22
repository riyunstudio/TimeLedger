package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secretKey := "mock-secret-key-for-development"

	claims := jwt.MapClaims{
		"user_type": "ADMIN",
		"user_id":   1,
		"center_id": 1,
		"exp":       time.Now().AddDate(0, 1, 0).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	fmt.Println(tokenString)
}
