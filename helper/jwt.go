package helper

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWTwithHMAC(tokenString string) (claims jwt.MapClaims, expired bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(" [F] invalid signing method ")
		}
		return []byte(os.Getenv("HUS_AUTH_TOKEN_KEY")), nil
	})
	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && (token.Valid || err == jwt.ErrTokenExpired) {
		if token.Valid {
			return claims, false, nil
		} else {
			return claims, true, nil
		}
	} else {
		return nil, false, fmt.Errorf(" [F] invalid token ")
	}
}
