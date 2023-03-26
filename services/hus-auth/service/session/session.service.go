package session

import (
	"context"
	"fmt"
	"hus-auth/common/hus"
	"hus-auth/db"
	"hus-auth/ent"
	"hus-auth/helper"

	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CreateHusSession takes user's uuid and create new hus session and return it.
func CreateHusSession(ctx context.Context, client *ent.Client, uid uuid.UUID, preserved bool) (
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
			"tid":     hs.Tid,                             // token id
			"purpose": "hus_session",                      // purpose
			"iss":     hus.AuthURL,                        // issuer
			"uid":     uid,                                // user's uuid
			"iat":     hs.Iat.Unix(),                      // issued at
			"exp":     time.Now().AddDate(0, 0, 7).Unix(), // expiration : one week
		})
	} else {
		rt = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sid":     hs.ID,                                // session token's uuid
			"tid":     hs.Tid,                               // token id
			"purpose": "hus_session",                        // purpose
			"iss":     hus.AuthURL,                          // issuer
			"uid":     uid,                                  // user's uuid
			"iat":     hs.Iat.Unix(),                        // issued at
			"exp":     time.Now().Add(time.Hour * 1).Unix(), // expiration : an hour
		})
	}

	hsk := []byte(hus.HusSecretKey)

	rts, err := rt.SignedString(hsk)
	if err != nil {
		return "", "", fmt.Errorf("!!signing hus-session token failed:%w", err)
	}
	log.Printf("hus-session created by (%s)", uid)
	// Sign and return the complete encoded token as a string
	return hs.ID.String(), rts, nil
}

func ValidateHusSession(ctx context.Context, client *ent.Client, hst string) (sid string, su *ent.User, err error) {
	claims, exp, err := helper.ParseJWTwithHMAC(hst)
	if err != nil {
		return "", nil, fmt.Errorf("invalid token")
	}

	hus_sid := claims["sid"].(string)
	hus_tid := claims["tid"].(string)
	hus_uid := claims["uid"].(string)

	if exp {
		return hus_sid, nil, fmt.Errorf("expired sesison")
	}
	// if the purpose is not hus_session, then return 401.
	if claims["purpose"].(string) != "hus_session" {
		return hus_sid, nil, fmt.Errorf("wrong purpose")
	}

	// check if the hus session is not revoked by querying the database with hus_sid.
	hs, err := db.QuerySessionBySID(ctx, client, hus_sid)
	if err != nil || hs == nil {
		return "", nil, fmt.Errorf("no such session")
	}
	/* for security if the token id is not matched, then revoke the session. */
	if hus_tid != hs.Tid.String() {
		return hus_sid, nil, fmt.Errorf("invalid token")
	}
	// check if the user exists by querying the database with hus_uid.
	u, err := db.QueryUserByUID(ctx, client, hus_uid)
	if err != nil || u == nil {
		return hus_sid, nil, fmt.Errorf("no such user")
	}
	return hus_sid, u, nil
}

func RefreshHusSession(ctx context.Context, client *ent.Client, sid string) (nstSigned string, err error) {
	sid_uuid, err := uuid.Parse(sid)
	if err != nil {
		return "", fmt.Errorf("invalid sid")
	}

	new_tid := uuid.New()
	hs, err := client.HusSession.UpdateOneID(sid_uuid).SetTid(new_tid).Save(ctx)
	if err != nil {
		return "", fmt.Errorf("updating session failed")
	}

	nst := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":     hs.ID.String(),                       // session token's uuid
		"tid":     hs.Tid.String(),                      // token id
		"purpose": "hus_session",                        // purpose
		"iss":     hus.AuthURL,                          // issuer
		"uid":     hs.UID,                               // user's uuid
		"iat":     hs.Iat.Unix(),                        // issued at
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // expiration : an hour
	})

	nstSigned, err = nst.SignedString([]byte(hus.HusSecretKey))
	if err != nil {
		return "", fmt.Errorf("signing hus_st failed")
	}

	return nstSigned, nil
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

// RevokeHusSessionToken takes session token and revokes them.
func RevokeHusSessionToken(ctx context.Context, client *ent.Client, st string) error {
	stClaims, _, err := helper.ParseJWTwithHMAC(st)

	sid_uuid, err := uuid.Parse(stClaims["sid"].(string))
	if err != nil {
		return err
	}

	err = client.HusSession.DeleteOneID(sid_uuid).Exec(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	return nil
}
