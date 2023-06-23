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

// V2

// HusSessionHandler godoc
// @Tags         auth
// @Router /hussession [get]
// @Summary checks and issues the Hus session token
// @Description this endpoint can be used both for Cloudhus and subservices.
// @Description if the subservice redirects the client to this endpoint with service name, session id and redirect url, its session will be connected to Hus session.
// @Description and if fallback url is given, it will redirect to fallback url if it fails.
// @Description but if they are not given, it will just respond rather than redirecting.
// @Description note that all urls must be url-encoded.
// @Param service query string false "subservice name"
// @Param sid query string false "subservice session id"
// @Param redirect query string false "redirect url"
// @Param fallback query string false "fallback url"
// @Success      200 "Ok, validated and connected"
// @Success      201 "Created, new Hus session connected"
// @Success      303 "See Other, redirection"
// @Failure	  400 "Bad Request"
// @Failure 500 "Internal Server Error"
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

	// if any of three parameters are given, either all or none of them must be given.
	if (serviceName != "" || sessionID != "" || redirectURL != "") && (serviceName == "" || sessionID == "" || redirectURL == "") {
		return c.String(http.StatusBadRequest, "service, sid, redirect should be given all together or none")
	}

	// if the request comes from subservice, handle the query parameters.
	var sessionUUID uuid.UUID
	if serviceName != "" {
		_, ok := common.Subservice[serviceName]
		// if the service name is not registered, return error.
		if !ok {
			return c.String(http.StatusBadRequest, "no such service")
		}
		// sessionID to UUID
		sessionUUID, err = uuid.Parse(sessionID)
		if err != nil {
			return c.String(http.StatusBadRequest, "sid is not valid")
		}
	}

	// url decode
	redirectURL, err1 := url.QueryUnescape(redirectURL)
	fallbackURL, err2 := url.QueryUnescape(fallbackURL)
	if err1 != nil || err2 != nil {
		return c.String(http.StatusBadRequest, "url is not valid")
	}

	// when there's no valid Hus session, create new one depending on this flag.
	createFlag := false

	// get hus_st from cookie
	hus_st, err := c.Cookie("hus_st")
	if err == http.ErrNoCookie || hus_st.Value == "" {
		// create new Hus session
		createFlag = true
	} else if err != nil {
		// there's an error while getting cookie, return error.
		if fallbackURL != "" {
			return c.Redirect(http.StatusSeeOther, fallbackURL)
		}
		return c.String(http.StatusInternalServerError, "an error occured while getting cookie")
	}

	// HUS SESSION EXSISTS, NOW VALIDATE IT

	// first validate and parse the session token
	var hs *ent.HusSession
	var preserved bool
	err = nil
	if !createFlag {
		hs, _, preserved, err = session.ValidateHusSessionV2(c.Request().Context(), ac.dbClient, hus_st.Value)
	}
	if err != nil || createFlag {
		/* NEW HUS SESSION CREATION */
		CreateHusSessionParams := session.CreateHusSessionParams{
			Ctx:     c.Request().Context(),
			Dbc:     ac.dbClient,
			Service: &serviceName,
			Sid:     &sessionUUID,
		}
		_, nhst, err := session.CreateHusSessionV2(CreateHusSessionParams)
		if err != nil {
			if fallbackURL != "" {
				return c.Redirect(http.StatusSeeOther, fallbackURL)
			}
			return c.String(http.StatusInternalServerError, "an error occured while creating hus session")
		}
		nhstCookie := &http.Cookie{
			Name:     "hus_st",
			Value:    nhst,
			Path:     "/",
			Domain:   hus.AuthCookieDomain,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		c.SetCookie(nhstCookie)

		if redirectURL != "" {
			return c.Redirect(http.StatusSeeOther, redirectURL)
		}
		return c.String(http.StatusCreated, "new Hus session created")
	}

	// HANDLE THE VALID SESSION

	// if the request comes from subservice, connect the sessions.
	if serviceName != "" {
		err = session.ConnectSessions(c.Request().Context(), ac.dbClient, hs, serviceName, sessionUUID)
		if err != nil {
			if fallbackURL != "" {
				return c.Redirect(http.StatusSeeOther, fallbackURL)
			}
			return c.String(http.StatusInternalServerError, "connecting sessions failed")
		}
	}

	// finally, rotate the Hus session
	nhst, err := session.RotateHusSessionV2(c.Request().Context(), ac.dbClient, hs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "rotating hus session failed")
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
	if preserved {
		nhstCookie.Expires = time.Now().Add(7 * 24 * time.Hour)
	}

	c.SetCookie(nhstCookie)

	if redirectURL != "" {
		return c.Redirect(http.StatusSeeOther, redirectURL)
	}
	return c.String(http.StatusOK, "valid Hus session")
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
