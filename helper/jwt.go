package helper

import (
	"context"
	"fmt"
	"hus-auth/ent"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateRefreshToken(c context.Context, client *ent.Client, token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return os.Getenv("HUS_AUTH_TOKEN_KEY"), nil
	})
	if err != nil {
		return nil, err
	}
	// claims exists
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid refresh token")
	}
}
