package helper

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWTwithHMAC(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(" [F] invalid signing method ")
		}
		return []byte(os.Getenv("HUS_AUTH_TOKEN_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf(" [F] invalid token ")
	}
}
