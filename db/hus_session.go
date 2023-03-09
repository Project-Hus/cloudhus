package db

import (
	"context"
	"hus-auth/ent"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CreateHusSession takes user's uuid and create new hus session and return it.
// and it also takes previous session ids as optional argument to revoke them.
func CreateHusSession(ctx context.Context, client *ent.Client, uid uuid.UUID, exp bool, prevSid ...string) (string, error) {
	// Delete previous session ids in prevSid
	for _, sid := range prevSid {
		sid, err := uuid.FromBytes([]byte(sid))
		if err != nil {
			log.Println("[F] converting sid to uuid failed:", err)
			return "", err
		}
		err = client.HusSession.DeleteOneID(sid).Exec(ctx)
		if err != nil {
			log.Print("[F] deleting previous session failed: ", err)
			return "", err
		}
	}

	// create new Hus session in database
	nhs := client.HusSession.Create().SetUID(uid)
	if exp { // if it's set to expired, give it 7 days expiration
		nhs = nhs.SetExp(time.Now().Add(time.Hour * 24 * 7))
	}
	hs, err := nhs.Save(ctx)
	if err != nil {
		log.Println("[F] creating new hus session failed:", err)
	}

	// using created id which is uuid, create refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tid":     hs.ID,                     // refresh token's uuid
		"purpose": "hus_session",             // purpose
		"iss":     os.Getenv("HUS_AUTH_URL"), // issuer
		"uid":     uid,                       // user's uuid
		"iat":     hs.Iat,                    // issued at
		"exp":     hs.Exp,                    // expiration : one week
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
