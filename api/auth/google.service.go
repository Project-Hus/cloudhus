package auth

import (
	"fmt"
	"hus-auth/db"
	"hus-auth/ent"
	"hus-auth/helper"
	"hus-auth/service/session"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
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
	// client ID that Google issued to lifthus
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	serviceParam := c.Param("service")
	var serviceUrl string
	switch serviceParam {
	case "lifthus":
		serviceUrl = os.Getenv("LIFTHUS_URL")
	}

	// credential sent from Google
	credential := c.FormValue("credential")

	// get where the user redirected from
	fmt.Println(c.FormValue("redirect"))

	// validating and parsing Google ID token
	payload, err := idtoken.Validate(c.Request().Context(), credential, clientID)
	if err != nil {
		log.Println("[F] Invalid ID token: %w", err)
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}
	// check if the user's ID token was intended for Hus.
	if payload.Audience != clientID {
		log.Println("[F] Invalid client ID:", payload.Audience)
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	// Google's unique user ID
	sub := payload.Claims["sub"].(string)
	// check if the user is registered with Google
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		_, err := db.CreateUserFromGoogle(c.Request().Context(), ac.Client, *payload)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
		}
	}

	// We checked or created if the Google user exists in Hus above,
	// Now get user query again to create new hus session.
	u, err = db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	_, HusSessionTokenSigned, err := session.CreateNewHusSession(c.Request().Context(), ac.Client, u.ID, false)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	}

	// set cookie for refresh token with 7 days expiration by struct literal
	cookie := &http.Cookie{
		Name:     "hus_st",
		Value:    HusSessionTokenSigned,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Domain:   os.Getenv("HUS_AUTH_DOMAIN"),
		SameSite: http.SameSiteDefaultMode,
	}
	c.SetCookie(cookie)

	// redirects to {serviceUrl}/hus/token/{hus-session-id}
	return c.Redirect(http.StatusMovedPermanently, serviceUrl)
}

// SessionRevocationHandler godoc
// @Router       /session/revoke [delete]
// @Summary      gets hus session token from cookie and revoke it.
// @Description  gets hus session token from cookie and revoke it by deleting it from database.
// @Tags         auth
// @Param        jwt header string true "Hus session token in cookie"
// @Success      200 "Ok"
// @Failure      500 "doesn't have to be handled"
func (ac authApiController) SessionRevocationHandler(c echo.Context) error {
	// get hus_st from cookie
	hus_st, _ := c.Cookie("hus_st")
	if hus_st == nil {
		return c.NoContent(http.StatusOK)
	}
	// Revoke past session in cookie

	claims, _, err := helper.ParseJWTwithHMAC(hus_st.Value)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	sid := claims["sid"].(string)

	suuid, err := uuid.Parse(sid)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	err = ac.Client.HusSession.DeleteOneID(suuid).Exec(c.Request().Context())
	if err != nil {
		if !ent.IsNotFound(err) {
			log.Print("[F] deleting past session failed: ", err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	// delete the session from cookie
	cookie := &http.Cookie{
		Name:     "hus_st",
		Value:    "",
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Domain:   os.Getenv("HUS_AUTH_DOMAIN"),
		SameSite: http.SameSiteDefaultMode,
	}
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}
