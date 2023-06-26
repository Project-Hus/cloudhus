package auth

import (
	"strconv"
	"strings"

	"fmt"
	"hus-auth/common"
	"hus-auth/common/helper"
	"hus-auth/common/hus"
	"hus-auth/common/service/session"
	"hus-auth/ent"
	"hus-auth/ent/connectedsession"

	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// HusSessionCheckHandler godoc
// @Router /session/check/{service}/{sid} [post]
// @Summary chekcs the service and sid and tells the subservice server that the client is signed in.
// @Description checks the service and sid and tells the subservice server that the client is signed in.
// @Description after the subservice server updates the session and responds with 200,
// @Description Hus auth server also reponds with 200 to the client.
// @Tags         auth
// @Param service path string true "subservice name"
// @Param sid path string true "session id"
// @Success      200 "Ok, theclient now should go to subservice's signing endpoint"
// @Failure      401 "Unauthorized, the client is not signed in"
// @Failure 404 "Not Found, the service is not registered"
// @Failure 500 "Internal Server Error"
func (ac authApiController) HusSessionCheckHandler(c echo.Context) error {
	// get service name and sid from path
	service := c.Param("service")
	lifthus_sid := c.Param("sid")

	subservice, ok := common.Subservice[service]
	// if the service name is not registered, return 404
	if !ok {
		return c.String(http.StatusNotFound, "no such service")
	}

	// get hus_st from cookie
	hus_st, err := c.Cookie("hus_st")
	// no valid session token, then return 401
	if err != nil || hus_st.Value == "" {
		return c.String(http.StatusUnauthorized, "not sigend in")
	}

	// first validate and parse the session token and get SID, User entity.
	hus_sid, u, err := session.ValidateHusSession(c.Request().Context(), ac.dbClient, hus_st.Value)
	if err != nil {
		sidUUID, _ := uuid.Parse(hus_sid)
		_ = ac.dbClient.HusSession.DeleteOneID(sidUUID).Exec(c.Request().Context())
		return c.String(http.StatusUnauthorized, "not signed in")
	}

	// if session token is valid, rotate the session token in cookie.
	nhstSigned, err := session.RefreshHusSession(c.Request().Context(), ac.dbClient, hus_sid)
	if err != nil {
		err = fmt.Errorf("refreshing hus session failed:%w", err)
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	nhstCookie := &http.Cookie{
		Name:     "hus_st",
		Value:    nhstSigned,
		Path:     "/",
		Domain:   hus.AuthCookieDomain,
		Expires:  time.Now().Add(time.Hour * 1),
		HttpOnly: true,
		Secure:   hus.CookieSecure,
		SameSite: hus.SameSiteMode,
	}
	c.SetCookie(nhstCookie)

	var bd string
	if u.Birthdate != nil {
		bd = u.Birthdate.Format(time.RFC3339)
	}

	hscbJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":               lifthus_sid,
		"uid":               strconv.FormatUint(u.ID, 10),
		"profile_image_url": u.ProfilePictureURL,
		"email":             u.Email,
		"email_verified":    u.EmailVerified,
		"name":              u.Name,
		"given_name":        u.GivenName,
		"family_name":       u.FamilyName,
		"birthdate":         bd,
		"exp":               time.Now().Add(time.Second * 10).Unix(),
	})

	hscbSigned, err := hscbJWT.SignedString([]byte(hus.HusSecretKey))
	if err != nil {
		err = fmt.Errorf("signing jwt for %s failed:%w", service, err)
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// with ac.httpClient, transfer the validation result to subservice auth server.
	req, err := http.NewRequest("PATCH", subservice.Subdomains["auth"].URL+"/auth/hus/session/sign", strings.NewReader(hscbSigned))
	if err != nil {
		err = fmt.Errorf("session injection to "+subservice.Domain.Name+" failed:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// send the request
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("session injection to "+subservice.Domain.Name+" failed:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return c.String(http.StatusOK, "session injection to "+subservice.Domain.Name+" success")
	} else {
		log.Println("an error occured from " + subservice.Domain.Name + ":" + resp.Status)
		return c.NoContent(http.StatusInternalServerError)
	}
}

// SessionRevocationHandler godoc
// @Router       /session/revoke [delete]
// @Summary      revokes every hus session in cookie from database.
// @Description  can be used to sign out.
// @Tags         auth
// @Param        jwt header string false "Hus session tokens in cookie"
// @Success      200 "Ok"
// @Failure      500 "doesn't have to be handled"
func (ac authApiController) SessionRevocationHandler(c echo.Context) error {
	stsToRevoke := []string{}

	// get hus_st from cookie
	hus_st, _ := c.Cookie("hus_st")
	hus_pst, _ := c.Cookie("hus_pst")
	if hus_st != nil && hus_st.Value != "" {
		stsToRevoke = append(stsToRevoke, hus_st.Value)
	}
	if hus_pst != nil && hus_pst.Value != "" {
		stsToRevoke = append(stsToRevoke, hus_pst.Value)
	}

	// Revoke all captured session tokens
	for _, st := range stsToRevoke {
		err := session.RevokeHusSessionToken(c.Request().Context(), ac.dbClient, st)
		if err != nil {
			err = fmt.Errorf("revoking hus session failed:%w", err)
			log.Println(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	cookie := &http.Cookie{
		Name:     "hus_st",
		Value:    "",
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Domain:   hus.AuthCookieDomain,
		SameSite: hus.SameSiteMode,
	}
	c.SetCookie(cookie)

	cookie2 := &http.Cookie{
		Name:     "hus_pst",
		Value:    "",
		Path:     "/",
		Secure:   hus.CookieSecure,
		HttpOnly: true,
		Domain:   hus.AuthCookieDomain,
		SameSite: hus.SameSiteMode,
	}
	c.SetCookie(cookie2)

	return c.NoContent(http.StatusOK)
}

// V2 ===================================================================================

// HusSessionHandler godoc
// @Tags         auth
// @Router /hussession [get]
// @Summary checks and issues the Hus session token
// @Description this endpoint can be used both for Cloudhus and subservices.
// @Description if the subservice redirects the client to this endpoint with service name, session id and redirect url, its session will be connected to Hus session.
// @Description and if fallback url is given, it will redirect to fallback url if it fails.
// @Description note that all urls must be url-encoded.
// @Param service query string true "subservice name"
// @Param redirect query string true "redirect url"
// @Param fallback query string false "fallback url"
// @Param sid query string true "subservice session id"
// @Success      303 "See Other, redirection"
// @Failure      303 "See Other, redirection"
func (ac authApiController) HusSessionHandler(c echo.Context) error {
	var err error

	// Query Parameters
	serviceName := c.QueryParam("service")  // name of the subservice that is requesting
	sessionID := c.QueryParam("sid")        // session ID of the subservice that is requesting
	redirectURL := c.QueryParam("redirect") // URL to be redirected after the request is processed
	fallbackURL := c.QueryParam("fallback") // URL to be redirected if the request fails
	if fallbackURL == "" {
		// if fallback URL is not given, it redirects to redirect URL
		fallbackURL = redirectURL
	}

	// if any of three parameters are not given, this request can't be handled.
	if serviceName == "" || sessionID == "" || redirectURL == "" {
		return c.Redirect(http.StatusSeeOther, common.Subservice["cloudhus"].Domain.URL+"/error")
	}

	// url decode
	redirectURL, err1 := url.QueryUnescape(redirectURL)
	fallbackURL, err2 := url.QueryUnescape(fallbackURL)
	if err1 != nil || err2 != nil {
		// invalid url
		return c.Redirect(http.StatusSeeOther, common.Subservice["cloudhus"].Domain.URL+"/error")
	}

	// service not registered, then halt.
	_, ok := common.Subservice[serviceName]
	if !ok {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}
	// sessionID to UUID
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	// when there's no valid Hus session, create new one depending on this flag.
	var hs *ent.HusSession
	createFlag := false

	// get hus_st from cookie
	hus_st, err := c.Cookie("hus_st")
	if err == http.ErrNoCookie || hus_st.Value == "" {
		// no session, create new one
		createFlag = true
	} else if err != nil {
		// there's an error while getting cookie, return error.
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	} else {
		hs, _, err = session.ValidateHusSessionV2(c.Request().Context(), hus_st.Value)
		if err != nil {
			createFlag = true
		}
	}

	// if no valid Hus session found, establish new Hus session.
	// after redirection to this same endpoint, it will handle newly established Hus session.
	if createFlag {
		/* NEW HUS SESSION CREATION */
		_, nhst, err := session.CreateHusSessionV2(c.Request().Context())
		if err != nil {
			return c.Redirect(http.StatusSeeOther, fallbackURL)
		}

		nhstCookie := &http.Cookie{
			Name:     "hus_st",
			Value:    nhst,
			Path:     "/",
			Domain:   hus.AuthCookieDomain,
			HttpOnly: true,
			Secure:   hus.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		}
		c.SetCookie(nhstCookie)

		// redirect to same endpoint here with same path and queries
		// this is to guarantee that the session is established between Cloudhus and the client
		tmpRedirect := common.Subservice["cloudhus"].Subdomains["auth"].URL +
			"/auth/hussession?service=" + serviceName +
			"&redirect=" + redirectURL +
			"&fallback=" + fallbackURL +
			"&sid=" + sessionID
		c.Redirect(http.StatusSeeOther, tmpRedirect)
	}

	// now handle the valid session
	err = session.ConnectSessions(c.Request().Context(), hs, serviceName, sessionUUID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	// finally, rotate the Hus session
	nhst, err := session.RotateHusSessionV2(c.Request().Context(), ac.dbClient, hs)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, fallbackURL)
	}

	// any kind of error(mostly Lambda timeout) may occur here after rotation before the user gets new tid.
	// this could be handled by user doing double check with another request.
	// or allowing the tid rotated only one step before. in this case new tid must be revoked.

	nhstCookie := &http.Cookie{
		Name:     "hus_st",
		Value:    nhst,
		Path:     "/",
		Domain:   hus.AuthCookieDomain,
		HttpOnly: true,
		Secure:   hus.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	if hs.Preserved {
		nhstCookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	}
	c.SetCookie(nhstCookie)

	return c.Redirect(http.StatusSeeOther, redirectURL)
}

// SessionConnectionHandler godoc
// @Tags         auth
// @Router /hussession/{token} [get]
// @Summary gets connection token from subservice and returns Hus session ID and user info
// @Description the token has properties pps, service and sid.
// @Param token path string true "pps, service name, session ID in signed token which expires only in 10 seconds"
// @Success      200 "Ok, session has been connected"
// @Failure      400 "Bad Request"
// @Failure      404 "Not Found, no such connected session"
func (ac authApiController) SessionConnectionHandler(c echo.Context) error {
	token := c.Param("token")
	claims, exp, err := helper.ParseJWTWithHMAC(token)
	if err != nil || exp {
		return c.String(http.StatusBadRequest, "invalid token")
	}

	pps := claims["pps"].(string)
	if pps != "session_connection" {
		return c.String(http.StatusBadRequest, "invalid token")
	}

	service := claims["service"].(string)
	sid := claims["sid"].(string)
	suuid, err := uuid.Parse(sid)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid token")
	}

	cs, err := ac.dbClient.ConnectedSession.Query().Where(connectedsession.And(
		connectedsession.Service(service),
		connectedsession.Csid(suuid),
	)).WithHusSession(func(hsq *ent.HusSessionQuery) {
		hsq.WithUser()
	}).Only(c.Request().Context())
	if err != nil {
		return c.String(http.StatusNotFound, "no such session")
	}

	return c.JSON(http.StatusOK, struct {
		Hsid string    `json:"hsid"`
		User *ent.User `json:"user,omitempty"`
	}{
		Hsid: cs.Hsid.String(),
		User: cs.Edges.HusSession.Edges.User,
	})
}

func (ac authApiController) SignOutHandler(c echo.Context) error {
	return nil
}
