package helper

import (
	"context"
	"fmt"
	"hus-auth/ent"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetNewAccessToken takes user's uuid and create signed access token and return it.
func GetNewAccessToken(c context.Context, client *ent.Client, uid string) (string, error) {
	aid := uuid.New().String()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aid":     aid,                       // refresh token's uuid
		"purpose": "access",                  // purpose
		"iss":     "https://api.lifthus.com", // issuer
		"uid":     uid,                       // user's uuid
		"iat":     time.Now().Unix(),         // issued at
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	accessTokenSigned, err := accessToken.SignedString([]byte(os.Getenv("HUS_AUTH_TOKEN_KEY")))
	if err != nil {
		log.Println("[F] signing access token failed: %w", err)
		return "", err
	}
	return accessTokenSigned, nil
}

// ValidateRefreshToken takes refresh token and validate it.
func ValidateRefreshToken(c context.Context, client *ent.Client, token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("HUS_AUTH_TOKEN_KEY")), nil
	})
	if err != nil {
		log.Println("[F] invalid refresh token:%w", err)
		return nil, err
	}
	// claims exists
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		log.Println("[F] invalid refresh token:%w", err)
		return nil, fmt.Errorf("invalid refresh token")
	}
}
