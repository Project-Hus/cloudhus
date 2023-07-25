package auth

import (
	"hus-auth/common"
	"hus-auth/common/db"
	"hus-auth/common/hus"
	"hus-auth/common/service/session"
	"hus-auth/ent/connectedsession"
	"hus-auth/ent/hussession"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

// GoogleAuthHandler godoc
// @Router       /hus/sign/social/google [post]
// @Summary      gets and processes Google ID token and redirects the user back to the given redirect url.
// @Description  validates the google ID token and do some authentication stuff.
// @Description  and redirects the user back to the given redirect url after the process is done.
// @Description  note that all urls must be url-encoded.
// @Tags         auth
// @Accept       json
// @Param redirect query string true "url to be redirected after authentication"
// @Param sid query string true "subservice session id"
// @Param fallback query string false "url to be redirected if the authentication fails"
// @Param        credential body string true "Google ID token"
// @Response      303 "See Other"
func (ac authApiController) GoogleAuthHandler(c echo.Context) error {
	// the session is already connected with subservice as the user accessed any page of the subservice.
	// so all this endpoint should do is just to validate the Google ID token and propagate the result to the connected sessions.

	redirectURL := c.QueryParam("redirect")
	csid := c.QueryParam("csid")

	if redirectURL == "" {
		return c.Redirect(http.StatusSeeOther, common.Subservice["cloudhus"].Domain.URL+"/error")
	}
	fallbackURL := c.QueryParam("fallback")
	if fallbackURL == "" {
		fallbackURL = redirectURL
	}

	redirectURL, err1 := url.QueryUnescape(redirectURL)
	fallbackURL, err2 := url.QueryUnescape(fallbackURL)
	if err1 != nil || err2 != nil {
		return c.Redirect(http.StatusSeeOther, common.Subservice["cloudhus"].Domain.URL+"/error")
	}

	csuuid, err := uuid.Parse(csid)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, common.Subservice["cloudhus"].Domain.URL+"/error")
	}

	/* Lax cookie not sent in cross-site redirect */
	// // validate the Hus session
	// hst, err := c.Cookie("hus_st")
	// if err != nil {
	// 	return c.Redirect(http.StatusSeeOther, fallbackURL)
	// }
	// hs, _, err := session.ValidateHusSession(c.Request().Context(), hst.Value)
	// if err != nil {
	// 	return c.Redirect(http.StatusSeeOther, fallbackURL)
	// }

	hs, err := db.Client.HusSession.Query().Where(hussession.HasConnectedSessionWith(connectedsession.Csid(csuuid))).
		WithConnectedSession().Only(c.Request().Context())
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
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), sub)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		u, err = db.CreateUserFromGoogle(c.Request().Context(), *payload)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, fallbackURL)
		}
	}

	err = session.SignHusSession(c.Request().Context(), hs, u)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	newToken, err := session.RotateHusSession(c.Request().Context(), hs)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	cookie := &http.Cookie{
		Name:     "hus_st",
		Value:    newToken,
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Domain:   hus.AuthCookieDomain,
		SameSite: hus.SameSiteMode,
	}
	if hs.Preserved {
		cookie.Expires = time.Now().AddDate(0, 0, 7)
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, redirectURL)
}
