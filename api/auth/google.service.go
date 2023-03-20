package auth

import (
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

	// revoke all previous hus sessions.
	stsToRevoke := []string{}
	hus_pst, _ := c.Cookie("hus_pst")
	hus_st, _ := c.Cookie("hus_st")
	if hus_pst != nil && hus_pst.Value != "" {
		stsToRevoke = append(stsToRevoke, hus_pst.Value)
	}
	if hus_st != nil && hus_st.Value != "" {
		stsToRevoke = append(stsToRevoke, hus_st.Value)
	}

	// revoke all captured hus session tokens
	for _, st := range stsToRevoke {
		err := session.RevokeHusSessionToken(c.Request().Context(), ac.dbClient, st)
		if err != nil {
			log.Println("revoking hus session token failed:", err)
			return c.String(http.StatusInternalServerError, "revoking hus session token failed")
		}
	}

	// client ID that Google issued to Cloudhus.
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

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
		log.Println("invalid id token:%w", err)
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
		Name:     "hus_pst",
		Value:    HusSessionTokenSigned,
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
