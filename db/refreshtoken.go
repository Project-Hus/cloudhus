package db

import (
	"context"
	"hus-auth/ent"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreateRefreshToken takes user's uuid and create signed refresh token and return it.
// It can be called only when user is authenticated by Third-party service.
func CreateRefreshToken(ctx context.Context, client *ent.Client, uid string) (string, error) {
	// first create refresh token in database, and its default key is uuid
	tk, err := client.RefreshToken.
		Create().SetUID(uid).Save(ctx)
	if err != nil {
		log.Print("[F] creating refresh token failed: ", err)
		return "", err
	}
	// using created id which is uuid, create refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tid":     tk.ID,                                     // refresh token's uuid
		"purpose": "refresh",                                 // purpose
		"iss":     os.Getenv("HOST_URL"),                     // issuer
		"uid":     uid,                                       // user's uuid
		"iat":     time.Now().Unix(),                         // issued at
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // expiration : one week
	})

	HATK := []byte(os.Getenv("HUS_AUTH_TOKEN_KEY"))

	rts, err := rt.SignedString(HATK)
	if err != nil {
		log.Print("[F] signing refresh token failed: ", err)
		return "", err
	}
	log.Printf("refresh token was created by (%s)", uid)
	// Sign and return the complete encoded token as a string
	return rts, nil
}
