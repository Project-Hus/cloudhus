package helper

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(" [F] invalid signing method ")
		}
		return []byte(os.Getenv("HUS_AUTH_TOKEN_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*jwt.MapClaims), nil
}
