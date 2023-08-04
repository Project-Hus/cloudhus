package session

import (
	"context"
	"fmt"
	"hus-auth/common"
	"hus-auth/common/db"
	"hus-auth/common/dto"
	"hus-auth/common/helper"
	"hus-auth/common/hus"
	"hus-auth/ent"
	"hus-auth/ent/hussession"
	"net/http"
	"strings"
	"sync"

	"time"

	"github.com/google/uuid"
)

func CreateHusSession(ctx context.Context) (newSession *ent.HusSession, newSignedToken string, err error) {
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

	hst, err := helper.SignedHST(hs)
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("generating session token failed:%w", err)
	}

	err = tx.Commit()
	if err != nil {
		err = db.Rollback(tx, err)
		return nil, "", fmt.Errorf("committing transaction failed:%w", err)
	}
	return hs, hst, nil
}

// ConnectSessions gets Hus session entity, subservice name and subservice's session ID.
// then makes connection between them and returns error if it occurs.
func ConnectSessions(ctx context.Context, hs *ent.HusSession, service string, csid uuid.UUID) error {
	if hs == nil {
		return fmt.Errorf("hussession is not given")
	}
	_, err := db.Client.ConnectedSession.Create().SetHsid(hs.ID).SetService(service).SetCsid(csid).Save(ctx)
	if err != nil {
		return fmt.Errorf("connecting sessions failed:%w", err)
	}
	return nil
}

// Session Error represents the error that occurs in session service package.
type SessionError struct {
	Message string
}

func (e SessionError) Error() string {
	return e.Message
}

// ExpiredSessionError occurs when the session token is expired.
var ExpiredValidSessionError = &SessionError{"expired valid session"}

// IsExpiredValid checks if the error is ExpiredValidSessionError.
func IsExpiredValid(err error) bool {
	return err == ExpiredValidSessionError
}

// ValidateHusSession gets Hus session token in string and validates it.
// if token is invalid, it returns "invalid session" error.
// if token is expired, it returns "expired session" error.
// if token's TID is not matched, it returns "illegal session" error after revoking user's all sessions.
// and if it is valid, it returns Hus session and User entities.
func ValidateHusSession(ctx context.Context, hst string) (hs *ent.HusSession, su *ent.User, err error) {
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
		return nil, nil, ExpiredValidSessionError
	}

	// check if the hus session is valid by querying the database with hus_sid.
	// and get the user entity too.
	hs, err = db.Client.HusSession.Query().Where(hussession.ID(husSid)).WithUser().Only(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid session")
	}

	su = hs.Edges.User

	// UUID type is a byte array with a length of 16.
	// so it can be compared directly.
	if hs.Tid != husTid {
		// revoke user's all session and propagate (not implemented yet) ------------------------------------------------------------------------
		if su != nil {
			_, _ = db.Client.HusSession.Delete().Where(hussession.UID(su.ID)).Exec(ctx)
		}
		return nil, nil, fmt.Errorf("illegal session")
	}

	return hs, su, nil
}

// RotateHusSession gets Hus session entity and rotates it's TID.
//
// any kind of error(mostly Lambda timeout) may occur after rotation before the user gets new tid.
// this could be handled by user doing double check with another request.
// or allowing the tid rotated only one step before. in this case new tid must be revoked.
func RotateHusSession(ctx context.Context, hs *ent.HusSession) (nstSigned string, err error) {
	hs, err = hs.Update().SetTid(uuid.New()).Save(ctx)
	if err != nil {
		return "", fmt.Errorf("updating session failed")
	}

	nst, err := helper.SignedHST(hs)
	if err != nil {
		return "", fmt.Errorf("signing Hus session failed")
	}

	return nst, nil
}

// SignHusSession takes Hus session entity and user entity and signs the Hus session.
// it also propagates to subservices which have connected session.
func SignHusSession(ctx context.Context, hs *ent.HusSession, u *ent.User) error {
	// sign the Hus session.
	hs, err := hs.Update().SetUID(u.ID).SetSignedAt(time.Now()).Save(ctx)
	if err != nil {
		return fmt.Errorf("signing Hus session failed:%w", err)
	}

	// query connected sessions if it is not queried yet.
	connectedSessions := hs.Edges.ConnectedSession
	if connectedSessions == nil {
		connectedSessions, err = hs.QueryConnectedSession().All(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("querying connected sessions failed:%w", err)
		}
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
			husConnectURL := service.Subdomains["auth"].URL + "/auth/hus/signin"

			husConnUser := &dto.HusConnUser{
				Uid:             u.ID,
				ProfileImageURL: u.ProfileImageURL,
				Email:           u.Email,
				EmailVerified:   u.EmailVerified,
				Name:            u.Name,
				GivenName:       u.GivenName,
				FamilyName:      u.FamilyName,
			}

			sipt, err := helper.SignedSIPToken(cs, *husConnUser)
			if err != nil {
				return
			}

			req, err := http.NewRequest(http.MethodPatch, husConnectURL, strings.NewReader(sipt))
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

// SignOutHusSession takes Hus session entity and signs out user's all Hus sessions.
func SignOutTotal(ctx context.Context, hsid uuid.UUID) error {
	// first query the hussession's owner and get all that user's hussessions with their connected sessions.
	hs, err := db.Client.HusSession.Query().Where(hussession.ID(hsid)).WithUser(func(q *ent.UserQuery) {
		q.WithHusSessions(func(q *ent.HusSessionQuery) {
			q.WithConnectedSession()
		})
	}).Only(ctx)
	if err != nil {
		return fmt.Errorf("querying hussession failed:%w", err)
	}

	if hs.Edges.User == nil {
		return fmt.Errorf("hussession is not signed")
	}
	userHusSessions := hs.Edges.User.Edges.HusSessions

	// gather hussessions' IDs.
	userHusSessionIDs := []uuid.UUID{}
	for _, hs := range hs.Edges.User.Edges.HusSessions {
		userHusSessionIDs = append(userHusSessionIDs, hs.ID)
	}
	// sign out user's all hussessions.
	err = db.Client.HusSession.Update().Where(hussession.IDIn(userHusSessionIDs...)).ClearUID().ClearSignedAt().Exec(ctx)
	if err != nil {
		return fmt.Errorf("total signing out failed:%w", err)
	}

	// gather all connected sessions.
	connectedSessions := []*ent.ConnectedSession{}
	for _, hs := range userHusSessions {
		connectedSessions = append(connectedSessions, hs.Edges.ConnectedSession...)
	}

	// asynchronically propagate to subservices.
	wg := sync.WaitGroup{}
	wg.Add(len(connectedSessions))
	for _, cs := range connectedSessions {
		go func(cs *ent.ConnectedSession) {
			defer wg.Done()
			service, ok := common.Subservice[cs.Service]
			if !ok {
				return
			}
			husConnectURL := service.Subdomains["auth"].URL + "/auth/hus/signout"

			sopt, err := helper.SignedSOPToken(cs)
			if err != nil {
				return
			}

			req, err := http.NewRequest(http.MethodPatch, husConnectURL, strings.NewReader(sopt))
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
