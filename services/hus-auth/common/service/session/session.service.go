package session

import (
	"context"
	"fmt"
	"hus-auth/common"
	"hus-auth/common/db"
	"hus-auth/common/helper"
	"hus-auth/common/hus"
	"hus-auth/ent"
	"hus-auth/ent/hussession"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CreateHusSessionV2 issues new Hus session and returns it.
func CreateHusSessionV2(ctx context.Context) (newSession *ent.HusSession, newSignedToken string, err error) {
	tx, err := db.Client.Tx(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("starting transaction failed:%w", err)
	}
	// create new Hus session
	hs, err := tx.HusSession.Create().Save(ctx)
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("creating new hus session failed:%w", err)
	}
	// Hus Session Token
	hst := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pps": "hus_session",
		"sid": hs.ID.String(),
		"tid": hs.Tid.String(),
		"iss": hus.AuthURL,
		"iat": hs.Iat.Unix(),
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"prv": hs.Preserved,
	})
	hsts, err := hst.SignedString(hus.HusSecretKeyBytes)
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("signing session token failed:%w", err)
	}

	err = tx.Commit()
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("committing transaction failed:%w", err)
	}
	return hs, hsts, nil
}

// ConnectSessions gets Hus session entity, subservice name and subservice's session ID.
// then makes connection between them and returns error if it occurs.
func ConnectSessions(ctx context.Context, hs *ent.HusSession, service string, csid uuid.UUID) error {
	_, err := db.Client.ConnectedSession.Create().SetHsid(hs.ID).SetService(service).SetCsid(csid).Save(ctx)
	if err != nil {
		return fmt.Errorf("connecting sessions failed:%w", err)
	}
	return nil
}

// ValidateHusSession gets Hus session token in string and validates it.
// if token is invalid, it returns "invalid session" error.
// if token is expired, it returns "expired session" error.
// if token's TID is not matched, it returns "illegal session" error after revoking all user's sessions.
// and if it is valid, it returns Hus session and User entities.
func ValidateHusSessionV2(ctx context.Context, hst string) (hs *ent.HusSession, su *ent.User, err error) {
	// parse the Hus session token.
	claims, exp, err := helper.ParseJWTWithHMAC(hst)
	if err != nil || claims["pps"].(string) != "hus_session" {
		return nil, nil, fmt.Errorf("invalid session")
	}
	// get and parse the Hus session ID and TID.
	husSidStr := claims["sid"].(string)
	husTidStr := claims["tid"].(string)
	husSid, err1 := uuid.Parse(husSidStr)
	husTid, err2 := uuid.Parse(husTidStr)
	if err1 != nil || err2 != nil {
		return nil, nil, fmt.Errorf("invalid session")
	}
	if exp {
		// revoke all related sessions (not implemented yet) ------------------------------------------------------------------------
		return nil, nil, fmt.Errorf("expired sesison")
	}

	// check if the hus session is valid by querying the database with hus_sid.
	// and get the user entity too.
	hs, err = db.Client.HusSession.Query().Where(hussession.ID(husSid)).WithUser().Only(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid session")
	}

	u := hs.Edges.User

	// UUID type is a byte array with a length of 16.
	// so it can be compared directly.
	if hs.Tid != husTid {
		// revoke all user's session and propagate (not implemented yet) ------------------------------------------------------------------------
		if u != nil {
			_, _ = db.Client.HusSession.Delete().Where(hussession.UID(u.ID)).Exec(ctx)
		}
		return nil, nil, fmt.Errorf("illegal session")
	}

	return hs, u, nil
}

// RotateHusSession gets Hus session entity and rotates it's TID.
//
// any kind of error(mostly Lambda timeout) may occur after rotation before the user gets new tid.
// this could be handled by user doing double check with another request.
// or allowing the tid rotated only one step before. in this case new tid must be revoked.
func RotateHusSessionV2(ctx context.Context, hs *ent.HusSession) (nstSigned string, err error) {
	hs, err = hs.Update().SetTid(uuid.New()).Save(ctx)
	if err != nil {
		return "", fmt.Errorf("updating session failed")
	}

	nst := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pps": "hus_session",   // purpose
		"sid": hs.ID.String(),  // session token's uuid
		"tid": hs.Tid.String(), // token id
		"iss": hus.AuthURL,     // issuer
		"iat": hs.Iat.Unix(),   // issued at
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"prv": hs.Preserved, // preserved
	})

	nstSigned, err = nst.SignedString(hus.HusSecretKeyBytes)
	if err != nil {
		return "", fmt.Errorf("signing Hus session failed")
	}

	return nstSigned, nil
}

// SignHusSession takes Hus session entity and user entity and signs the Hus session.
// it also propagates to subservices which have connected session.
func SignHusSession(ctx context.Context, hs *ent.HusSession, u *ent.User) error {
	// sign the Hus session.
	hs, err := hs.Update().SetUID(u.ID).SetSignedAt(time.Now()).Save(ctx)
	if err != nil {
		return fmt.Errorf("signing Hus session failed:%w", err)
	}

	// get connected sessions.
	connectedSessions, err := hs.QueryConnectedSession().All(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("querying connected sessions failed:%w", err)
	}

	// propagate to connected sessions.
	// some of them may fail. in this case subservice checks the session at next refresh and etc.
	wg := sync.WaitGroup{}
	wg.Add(len(connectedSessions))
	for _, cs := range connectedSessions {
		go func(cs *ent.ConnectedSession) {
			defer wg.Done()
			service, ok := common.Subservice[cs.Service]
			if !ok {
				return
			}
			husConnectURL := service.Subdomains["auth"].URL + "/auth/hussession/connect"

			hscJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"pps":               "session_sign",
				"hsid":              cs.Hsid.String(),
				"sid":               cs.Csid.String(),
				"uid":               strconv.FormatUint(u.ID, 10),
				"profile_image_url": u.ProfileImageURL,
				"email":             u.Email,
				"email_verified":    u.EmailVerified,
				"name":              u.Name,
				"given_name":        u.GivenName,
				"family_name":       u.FamilyName,
				"birthdate":         nil,
				"exp":               time.Now().Add(time.Second * 10).Unix(),
			})

			hscSigned, err := hscJWT.SignedString(hus.HusSecretKeyBytes)
			if err != nil {
				return
			}

			req, err := http.NewRequest(http.MethodPatch, husConnectURL, strings.NewReader(hscSigned))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "text/plain")
			resp, err := hus.Http.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
		}(cs)
	}

	wg.Wait()

	return nil
}

func SignOutTotal(ctx context.Context, hsid uuid.UUID) error {
	return nil
}

// ==========================================================================================

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
	claims, exp, err := helper.ParseJWTWithHMAC(hst)
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
	u, err := db.QueryUserByUID(ctx, hus_uid_uint64)
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
	stClaims, _, err := helper.ParseJWTWithHMAC(st)

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
