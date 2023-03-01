package db

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

// CreateRefreshToken takes user's uuid and create signed refresh token and return it.
// It can be called only when user is authenticated by Third-party service.
func CreateRefreshToken(ctx context.Context, client *ent.Client, uid string) (string, error) {
	// UUID for refresh token
	tid := uuid.New().String()
	// Create new refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tid": tid,                                       // refresh token's uuid
		"iss": "https://api.lifthus.com",                 // issuer
		"uid": uid,                                       // user's uuid
		"iat": time.Now().Unix(),                         // issued at
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // expiration : one week
	})
	// Save refresh token to database
	_, err := client.RefreshToken.
		Create().SetID(tid).SetUID(uid).Save(ctx)
	if err != nil {
		return "", fmt.Errorf("failed creating refresh token: %w", err)
	}
	log.Println("refresh token was created: ", rt)
	// Sign and return the complete encoded token as a string
	return rt.SignedString(os.Getenv("HUS_AUTH_TOKEN_KEY"))
}
