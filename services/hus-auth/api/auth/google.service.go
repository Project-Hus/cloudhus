package auth

import (
	"hus-auth/common"
	"hus-auth/common/db"
	"hus-auth/common/hus"
	"hus-auth/common/service/session"
	"log"
	"net/http"
	"net/url"
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
// @Param subservice_name path string true "subservice name"
// @Param        jwt body string true "Google ID token"
// @Response      301 "to /auth/{token_string} or to /error"
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

	cookie := &http.Cookie{
		Name:     "hus_st",
		Value:    HusSessionTokenSigned,
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Domain:   hus.AuthCookieDomain,
		SameSite: hus.SameSiteMode,
	}
	c.SetCookie(cookie)

	cookie2 := &http.Cookie{
		Name:     "hus_pst",
		Value:    HusSessionTokenSigned,
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Expires:  time.Now().AddDate(1, 0, 0),
		Domain:   hus.AuthCookieDomain,
		SameSite: hus.SameSiteMode,
	}
	c.SetCookie(cookie2)

	cookieTest := &http.Cookie{
		Name:     "hus_test",
		Value:    "TESTCOOKIEHAPPYCOOKIE",
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Domain:   hus.AuthCookieDomain,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookieTest)

	cookieTest2 := &http.Cookie{
		Name:     "hus_test2",
		Value:    "TESTCOOKIEHAPPYCOOKIE",
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Domain:   hus.AuthCookieDomain,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookieTest2)

	// redirects to {serviceUrl}/hus/token/{hus-session-id}
	return c.Redirect(http.StatusMovedPermanently, serviceUrl)
}

// GoogleAuthHandlerV2 godoc
// @Router       /social/google [post]
// @Summary      gets and processes Google ID token and redirects the user back to the given redirect url.
// @Description  validates the google ID token and do some authentication stuff.
// @Description  and redirects the user back to the given  after the process.
// @Description  note that all urls must be url-encoded.
// @Tags         auth
// @Accept       json
// @Param redirect query string true "url to be redirected after authentication"
// @Param fallback query string false "url to be redirected if the authentication fails"
// @Param        credential body string true "Google ID token"
// @Response      303 "See Other"
func (ac authApiController) GoogleAuthHandlerV2(c echo.Context) error {
	// the session is already connected with subservice as the user accessed any page of the subservice.
	// so all this endpoint should do is just to validate the Google ID token and propagate the result to the connected sessions.

	redirectURL := c.QueryParam("redirect")
	fallbackURL := c.QueryParam("fallback")
	if fallbackURL == "" {
		fallbackURL = redirectURL
	}

	redirectURL, err1 := url.QueryUnescape(redirectURL)
	fallbackURL, err2 := url.QueryUnescape(fallbackURL)
	if err1 != nil || err2 != nil {
		return c.Redirect(http.StatusSeeOther, common.Subservice["cloudhus"].Subdomains["auth"].URL+"/auth")
	}

	if redirectURL == "" {
		redirectURL = common.Subservice["cloudhus"].Subdomains["auth"].URL + "/auth"
		fallbackURL = redirectURL
	}

	// validate the Hus session
	hst, err := c.Cookie("hus_st")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}
	hs, _, preserved, err := session.ValidateHusSessionV2(c.Request().Context(), ac.dbClient, hst.Value)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	// validate and parse the Google ID token
	payload, err := idtoken.Validate(c.Request().Context(), c.FormValue("credential"), hus.GoogleClientID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	// Google's unique user ID
	sub := payload.Claims["sub"].(string)
	// check if the user is registered with Google)
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.dbClient, sub)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		_, err := db.CreateUserFromGoogle(c.Request().Context(), ac.dbClient, *payload)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, fallbackURL)
		}
	}

	// // We checked or created if the Google user exists in Hus above,
	// // Now get user query again to create new hus session.
	// u, err = db.QueryUserByGoogleSub(c.Request().Context(), ac.dbClient, sub)
	// if err != nil {
	// 	return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	// }

	// _, HusSessionTokenSigned, err := session.CreateHusSession(c.Request().Context(), ac.dbClient, u.ID, false)
	// if err != nil {
	// 	return c.Redirect(http.StatusMovedPermanently, serviceUrl+"/error")
	// }

	// cookie := &http.Cookie{
	// 	Name:     "hus_st",
	// 	Value:    HusSessionTokenSigned,
	// 	Path:     "/",
	// 	Secure:   hus.CookieSecure,
	// 	HttpOnly: true,
	// 	Domain:   hus.AuthCookieDomain,
	// 	SameSite: hus.SameSiteMode,
	// }
	// c.SetCookie(cookie)

	// cookie2 := &http.Cookie{
	// 	Name:     "hus_pst",
	// 	Value:    HusSessionTokenSigned,
	// 	Path:     "/",
	// 	Secure:   hus.CookieSecure,
	// 	HttpOnly: true,
	// 	Expires:  time.Now().AddDate(1, 0, 0),
	// 	Domain:   hus.AuthCookieDomain,
	// 	SameSite: hus.SameSiteMode,
	// }
	// c.SetCookie(cookie2)

	// cookieTest := &http.Cookie{
	// 	Name:     "hus_test",
	// 	Value:    "TESTCOOKIEHAPPYCOOKIE",
	// 	Path:     "/",
	// 	Secure:   hus.CookieSecure,
	// 	HttpOnly: true,
	// 	Domain:   hus.AuthCookieDomain,
	// 	SameSite: http.SameSiteLaxMode,
	// }
	// c.SetCookie(cookieTest)

	// cookieTest2 := &http.Cookie{
	// 	Name:     "hus_test2",
	// 	Value:    "TESTCOOKIEHAPPYCOOKIE",
	// 	Path:     "/",
	// 	Secure:   hus.CookieSecure,
	// 	HttpOnly: true,
	// 	Domain:   hus.AuthCookieDomain,
	// 	SameSite: http.SameSiteStrictMode,
	// }
	// c.SetCookie(cookieTest2)

	// // redirects to {serviceUrl}/hus/token/{hus-session-id}
	// return c.Redirect(http.StatusMovedPermanently, serviceUrl)
}
