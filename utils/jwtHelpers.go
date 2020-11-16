package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"internship_project/models"
	"time"
)

var hmacSampleSecret []byte = []byte("my_secret_key")

func ParseJWT(jwt_string string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwt_string, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func CreateJWT(u models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"name": u.Name,
		"iat": time.Now().Unix(), // Issued at
		"exp": time.Now().Add(time.Hour * 2).Unix(), // Expires
		"nbf": time.Now().Add(time.Second * 5).Unix(), // Not before
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

