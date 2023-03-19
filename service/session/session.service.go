package session

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

// CreateHusSession takes user's uuid and create new hus session and return it.
func CreateNewHusSession(ctx context.Context, client *ent.Client, uid uuid.UUID, preserved bool) (
	new_sid, new_token string, err error,
) {
	// create new Hus session in database
	hs, err := client.HusSession.Create().SetUID(uid).SetPreserved(preserved).Save(ctx)
	if err != nil {
		return "", "", fmt.Errorf("!!creating new hus session failed:%w", err)
	}

	var rt *jwt.Token

	// using created session's UUID, create session token
	if preserved {
		rt = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sid":     hs.ID,                              // session token's uuid
			"purpose": "hus_session",                      // purpose
			"iss":     os.Getenv("HUS_AUTH_URL"),          // issuer
			"uid":     uid,                                // user's uuid
			"iat":     hs.Iat.Unix(),                      // issued at
			"exp":     time.Now().AddDate(0, 0, 7).Unix(), // expiration : one week
		})
	} else {
		rt = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sid":     hs.ID,                                // session token's uuid
			"purpose": "hus_session",                        // purpose
			"iss":     os.Getenv("HUS_AUTH_URL"),            // issuer
			"uid":     uid,                                  // user's uuid
			"iat":     hs.Iat.Unix(),                        // issued at
			"exp":     time.Now().Add(time.Hour * 1).Unix(), // expiration : an hour
		})
	}

	hsk := []byte(os.Getenv("HUS_SECRET_KEY"))

	rts, err := rt.SignedString(hsk)
	if err != nil {
		return "", "", fmt.Errorf("!!signing hus-session token failed:%w", err)
	}
	log.Printf("hus-session created by (%s)", uid)
	// Sign and return the complete encoded token as a string
	return hs.ID.String(), rts, nil
}

// RevokeHusSession takes session id and revokes it.
func RevokeHusSession(ctx context.Context, client *ent.Client, sid string) error {
	sid_uuid, err := uuid.Parse(sid)
	if err != nil {
		return fmt.Errorf("!!invalid sid:%w", err)
	}
	err = client.HusSession.DeleteOneID(sid_uuid).Exec(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("!!revoking hus session failed:%w", err)
	}
	return nil
}
