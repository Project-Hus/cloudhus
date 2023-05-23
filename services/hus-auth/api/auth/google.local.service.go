package auth

import (
	"hus-auth/common"
	"hus-auth/common/hus"
	"hus-auth/db"
	"hus-auth/service/session"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

// googleAuthHandler is a local development version of GoogleAuthHandler.
// which uses Authorization header instead of cookie.
func (ac authApiController) googleAuthHandler(c echo.Context) error {
	// get Authorization header
	authorization := c.Request().Header.Get("Authorization")
	stsToRevoke := []string{}
	stsToRevoke = append(stsToRevoke, authorization)

	// revoke all captured hus session tokens
	for _, st := range stsToRevoke {
		err := session.RevokeHusSessionToken(c.Request().Context(), ac.dbClient, st)
		if err != nil {
			log.Println("revoking hus session token failed:", err)
			return c.String(http.StatusInternalServerError, "revoking hus session token failed")
		}
	}

	// client ID that Google issued to Cloudhus.
	clientID := hus.GoogleClientID

	// check if the service is registered.
	serviceParam := c.Param("service")
	subservice, ok := common.Subservice[serviceParam]
	if !ok {
		return c.String(http.StatusNotFound, "no such service")
	}
	serviceUrl := subservice.Domain.URL

	// credential sent from Google
	credential := c.FormValue("credential")

	// validate and parse the Google ID token
	payload, err := idtoken.Validate(c.Request().Context(), credential, clientID)
	if err != nil {
		log.Println("invalid id token:", err)
		log.Println("@credential:", credential)
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	// Google's unique user ID
	sub := payload.Claims["sub"].(string)
	// check if the user is registered with Google
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.dbClient, sub)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		_, err := db.CreateUserFromGoogle(c.Request().Context(), ac.dbClient, *payload)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
		}
	}

	// We checked or created if the Google user exists in Hus above,
	// Now get user query again to create new hus session.
	u, err = db.QueryUserByGoogleSub(c.Request().Context(), ac.dbClient, sub)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	_, HusSessionTokenSigned, err := session.CreateHusSession(c.Request().Context(), ac.dbClient, u.ID, false)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	// set Authorization header with HusSessionTokenSigned
	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+HusSessionTokenSigned)

	// redirects to {serviceUrl}/hus/token/{hus-session-id}
	return c.Redirect(http.StatusMovedPermanently, serviceUrl)
}
