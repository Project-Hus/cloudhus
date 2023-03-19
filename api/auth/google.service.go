package auth

import (
	"errors"
	"hus-auth/common"
	"hus-auth/db"
	"hus-auth/service/session"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

// GoogleAuthHandler godoc
// @Router       /social/google/{subservice_name} [post]
// @Summary      gets google IDtoken and redirect with hus session cookie.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the refresh token will be expired in 7 days.
// @Tags         auth
// @Accept       json
// @Param        jwt body string true "Google ID token"
// @Success      301 "to /auth/{token_string}"
// @Failure      301 "to /error"
func (ac authApiController) GoogleAuthHandler(c echo.Context) error {
	// from cookie get hus_pst
	hus_psid, err := c.Cookie("hus_psid")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Println("getting hus_psid cookie failed:", err)
		return c.String(http.StatusInternalServerError, "getting hus_psid cookie faield")
	}
	// if there is existing hus session, revoke it.
	if hus_psid != nil && hus_psid.Value != "" {
		// revoke the old session.
		err = session.RevokeHusSession(c.Request().Context(), ac.dbClient, hus_psid.Value)
		if err != nil {
			log.Println("rovoking hus session failed", err)
			return c.String(http.StatusInternalServerError, "revoking hus session failed")
		}
	}

	// client ID that Google issued to lifthus
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	serviceParam := c.Param("service")
	subservice, ok := common.Subservice[serviceParam]
	if !ok {
		return c.String(http.StatusNotFound, "no such service")
	}
	serviceUrl := subservice.Domain.URL

	// credential sent from Google
	credential := c.FormValue("credential")

	// validating and parsing Google ID token
	payload, err := idtoken.Validate(c.Request().Context(), credential, clientID)
	if err != nil {
		log.Println("invalid id token:%w", err)
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}
	// check if the user's ID token was intended for Hus.
	if payload.Audience != clientID {
		log.Println("invalid client id:", payload.Audience)
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

	nsid, HusSessionTokenSigned, err := session.CreateNewHusSession(c.Request().Context(), ac.dbClient, u.ID, false)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	cookie := &http.Cookie{
		Name:     "hus_st",
		Value:    HusSessionTokenSigned,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Domain:   os.Getenv("HUS_AUTH_DOMAIN"),
		SameSite: http.SameSiteDefaultMode,
	}
	c.SetCookie(cookie)

	cookie2 := &http.Cookie{
		Name:     "hus_psid",
		Value:    nsid,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().AddDate(1, 0, 0),
		Domain:   os.Getenv("HUS_AUTH_DOMAIN"),
		SameSite: http.SameSiteDefaultMode,
	}
	c.SetCookie(cookie2)

	// redirects to {serviceUrl}/hus/token/{hus-session-id}
	return c.Redirect(http.StatusMovedPermanently, serviceUrl)
}
