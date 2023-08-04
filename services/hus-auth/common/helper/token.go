package helper

import (
	"hus-auth/common/dto"
	"hus-auth/common/hus"
	"hus-auth/ent"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignedHST(hs *ent.HusSession) (string, error) {
	token := NewJWT(jwt.MapClaims{
		"pps": "hus_session",
		"sid": hs.ID.String(),
		"tid": hs.Tid.String(),
		"iss": hus.AuthURL,
		"iat": hs.Iat.Unix(),
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"prv": hs.Preserved,
	})
	return SignJWT(token)
}

func SignedSIPToken(cs *ent.ConnectedSession, hcu dto.HusConnUser) (string, error) {
	token := NewJWT(jwt.MapClaims{
		"pps":  "signin_propagation",
		"hsid": cs.Hsid.String(),
		"csid": cs.Csid.String(),
		"user": hcu,
		"exp":  time.Now().Add(time.Second * 10).Unix(),
	})
	return SignJWT(token)
}

func SignedSOPToken(cs *ent.ConnectedSession) (string, error) {
	token := NewJWT(jwt.MapClaims{
		"pps":  "signout_propagation",
		"hsid": cs.Hsid.String(),
		"exp":  time.Now().Add(time.Second * 10).Unix(),
	})
	return SignJWT(token)
}
