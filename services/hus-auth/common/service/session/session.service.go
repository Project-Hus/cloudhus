package session

import (
	"context"
	"fmt"
	"hus-auth/common"
	"hus-auth/common/db"
	"hus-auth/common/hus"
	"hus-auth/ent"
	"hus-auth/helper"
	"strconv"

	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CreateHusSessionParams struct {
	Ctx context.Context
	// if the request is from Cloudhus, below fields are not required.
	Dbc       *ent.Client
	Service   *common.ServiceDomain
	Propagate *string
	Sid       *uuid.UUID
}

// CreateHusSessionV2 issues new Hus session and returns it.
// and if subservice's session ID is provided, it will be connected to the Hus session.
// after the connection is established, subservice must verify it by asking to Cloudhus.
func CreateHusSessionV2(ps CreateHusSessionParams) (
	newSession *ent.HusSession, newToken string, err error,
) {
	tx, err := ps.Dbc.Tx(ps.Ctx)
	if err != nil {
		return nil, "", fmt.Errorf("starting transaction failed:%w", err)
	}

	// create new Hus session
	hs, err := tx.HusSession.Create().Save(ps.Ctx)
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("!!creating new hus session failed:%w", err)
	}

	// Hus Session Token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pps": "hus_session",
		"sid": hs.ID,
		"tid": hs.Tid,
		"iss": hus.AuthURL,
		"iat": hs.Iat.Unix(),
		"prv": hs.Preserved,
	})
	rts, err := rt.SignedString(hus.HusSecretKeyBytes)
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("signing session token failed:%w", err)
	}

	err = tx.Commit()
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("committing transaction failed:%w", err)
	}

	// if there are service and sid, connect them to the Huse session created above.
	if ps.Service != nil && ps.Sid != nil {
		tx, err = ps.Dbc.Tx(ps.Ctx)
		if err != nil {
			return nil, "", fmt.Errorf("starting transaction failed:%w", err)
		}

		_, err := tx.ConnectedSession.Create().SetHsid(hs.ID).SetService(ps.Service.Domain.Name).SetCsid(*ps.Sid).Save(ps.Ctx)
		if err != nil {
			err = db.Rollback(tx, err)
			return nil, "", fmt.Errorf("connecting sessions failed:%w", err)
		}

		if err != nil {
			err = db.Rollback(tx, err)
			return nil, "", fmt.Errorf("signing session token failed:%w", err)
		}
	}

	return hs, rts, nil
}

// CreateHusSession takes user's uuid and create new hus session and return it.
func CreateHusSession(ctx context.Context, client *ent.Client, uid uint64, preserved bool) (
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
			"uid":     strconv.FormatUint(uid, 10),        // user's uuid
			"iat":     hs.Iat.Unix(),                      // issued at
			"exp":     time.Now().AddDate(0, 0, 7).Unix(), // expiration : one week
		})
	} else {
		rt = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sid":     hs.ID,                                // session token's uuid
			"tid":     hs.Tid,                               // token id
			"purpose": "hus_session",                        // purpose
			"iss":     hus.AuthURL,                          // issuer
			"uid":     strconv.FormatUint(uid, 10),          // user's uuid
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
	hus_uid_uint64, err := strconv.ParseUint(hus_uid, 10, 64)
	if err != nil {
		return hus_sid, nil, fmt.Errorf("invalid uid")
	}
	u, err := db.QueryUserByUID(ctx, client, hus_uid_uint64)
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
		"uid":     strconv.FormatUint(*hs.UID, 10),      // user's uuid
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

// UUID version
// func CreateHusSession(ctx context.Context, client *ent.Client, uid uuid.UUID, preserved bool) (
// 	new_sid, new_token string, err error,
// ) {
// 	// create new Hus session in database
// 	hs, err := client.HusSession.Create().SetUID(uid).SetPreserved(preserved).Save(ctx)
// 	if err != nil {
// 		return "", "", fmt.Errorf("!!creating new hus session failed:%w", err)
// 	}

// 	var rt *jwt.Token

// 	// using created session's UUID, create session token
// 	if preserved {
// 		rt = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"sid":     hs.ID,                              // session token's uuid
// 			"tid":     hs.Tid,                             // token id
// 			"purpose": "hus_session",                      // purpose
// 			"iss":     hus.AuthURL,                        // issuer
// 			"uid":     uid,                                // user's uuid
// 			"iat":     hs.Iat.Unix(),                      // issued at
// 			"exp":     time.Now().AddDate(0, 0, 7).Unix(), // expiration : one week
// 		})
// 	} else {
// 		rt = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 			"sid":     hs.ID,                                // session token's uuid
// 			"tid":     hs.Tid,                               // token id
// 			"purpose": "hus_session",                        // purpose
// 			"iss":     hus.AuthURL,                          // issuer
// 			"uid":     uid,                                  // user's uuid
// 			"iat":     hs.Iat.Unix(),                        // issued at
// 			"exp":     time.Now().Add(time.Hour * 1).Unix(), // expiration : an hour
// 		})
// 	}

// 	hsk := []byte(hus.HusSecretKey)

// 	rts, err := rt.SignedString(hsk)
// 	if err != nil {
// 		return "", "", fmt.Errorf("!!signing hus-session token failed:%w", err)
// 	}
// 	log.Printf("hus-session created by (%s)", uid)
// 	// Sign and return the complete encoded token as a string
// 	return hs.ID.String(), rts, nil
// }

// func ValidateHusSession(ctx context.Context, client *ent.Client, hst string) (sid string, su *ent.User, err error) {
// 	claims, exp, err := helper.ParseJWTwithHMAC(hst)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("invalid token")
// 	}

// 	hus_sid := claims["sid"].(string)
// 	hus_tid := claims["tid"].(string)
// 	hus_uid := claims["uid"].(string)

// 	if exp {
// 		return hus_sid, nil, fmt.Errorf("expired sesison")
// 	}
// 	// if the purpose is not hus_session, then return 401.
// 	if claims["purpose"].(string) != "hus_session" {
// 		return hus_sid, nil, fmt.Errorf("wrong purpose")
// 	}

// 	// check if the hus session is not revoked by querying the database with hus_sid.
// 	hs, err := db.QuerySessionBySID(ctx, client, hus_sid)
// 	if err != nil || hs == nil {
// 		return "", nil, fmt.Errorf("no such session")
// 	}
// 	/* for security if the token id is not matched, then revoke the session. */
// 	if hus_tid != hs.Tid.String() {
// 		return hus_sid, nil, fmt.Errorf("invalid token")
// 	}
// 	// check if the user exists by querying the database with hus_uid.
// 	u, err := db.QueryUserByUID(ctx, client, hus_uid)
// 	if err != nil || u == nil {
// 		return hus_sid, nil, fmt.Errorf("no such user")
// 	}
// 	return hus_sid, u, nil
// }

// func RefreshHusSession(ctx context.Context, client *ent.Client, sid string) (nstSigned string, err error) {
// 	sid_uuid, err := uuid.Parse(sid)
// 	if err != nil {
// 		return "", fmt.Errorf("invalid sid")
// 	}

// 	new_tid := uuid.New()
// 	hs, err := client.HusSession.UpdateOneID(sid_uuid).SetTid(new_tid).Save(ctx)
// 	if err != nil {
// 		return "", fmt.Errorf("updating session failed")
// 	}

// 	nst := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sid":     hs.ID.String(),                       // session token's uuid
// 		"tid":     hs.Tid.String(),                      // token id
// 		"purpose": "hus_session",                        // purpose
// 		"iss":     hus.AuthURL,                          // issuer
// 		"uid":     hs.UID,                               // user's uuid
// 		"iat":     hs.Iat.Unix(),                        // issued at
// 		"exp":     time.Now().Add(time.Hour * 1).Unix(), // expiration : an hour
// 	})

// 	nstSigned, err = nst.SignedString([]byte(hus.HusSecretKey))
// 	if err != nil {
// 		return "", fmt.Errorf("signing hus_st failed")
// 	}

// 	return nstSigned, nil
// }

// // RevokeHusSession takes session id and revokes it.
// func RevokeHusSession(ctx context.Context, client *ent.Client, sid string) error {
// 	sid_uuid, err := uuid.Parse(sid)
// 	if err != nil {
// 		return fmt.Errorf("!!invalid sid:%w", err)
// 	}
// 	err = client.HusSession.DeleteOneID(sid_uuid).Exec(ctx)
// 	if err != nil && !ent.IsNotFound(err) {
// 		return fmt.Errorf("!!revoking hus session failed:%w", err)
// 	}
// 	return nil
// }

// // RevokeHusSessionToken takes session token and revokes them.
// func RevokeHusSessionToken(ctx context.Context, client *ent.Client, st string) error {
// 	stClaims, _, err := helper.ParseJWTwithHMAC(st)

// 	sid_uuid, err := uuid.Parse(stClaims["sid"].(string))
// 	if err != nil {
// 		return err
// 	}

// 	err = client.HusSession.DeleteOneID(sid_uuid).Exec(ctx)
// 	if err != nil && !ent.IsNotFound(err) {
// 		return err
// 	}
// 	return nil
// }
