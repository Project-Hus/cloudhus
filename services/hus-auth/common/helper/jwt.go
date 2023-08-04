package helper

import (
	"errors"
	"fmt"
	"hus-auth/common/hus"

	"github.com/golang-jwt/jwt/v5"
)

// ParseJWTWithHMAC parses any HS256 JWT token with HusSecretKey and returns claims.
// it doesn't return error when the token is expired.
// instead, it returns expired flag with other return values(claims, expired, error).
func ParseJWTWithHMAC(tokenString string) (claims jwt.MapClaims, expired bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(hus.HusSecretKey), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if token.Valid {
			return claims, false, nil
		} else {
			return claims, true, nil
		}
	} else {
		return nil, false, fmt.Errorf("invalid token")
	}
}

// NewJWT creates a new jwt token not signed.
func NewJWT(c jwt.MapClaims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c)
}

// SignJWT signs a jwt token with HusSecretKey.
func SignJWT(t *jwt.Token) (string, error) {
	return t.SignedString(hus.HusSecretKeyBytes)
}
