package auth

import (
	"strconv"
	"strings"

	"fmt"
	"hus-auth/common"
	"hus-auth/common/hus"
	"hus-auth/common/service/session"

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

// HusSessionHandler godoc
// @Tags         auth
// @Router /session [get]
// @Summary checks and issues the Hus session token
// @Description this endpoint can be used both for Cloudhus and subservices.
// @Description if the subservice redirects the client to this endpoint with service name, session id and redirect url, the session will be connected with Cloudhus session.
// @Description but if they are given, it will just respond rather than redirecting.
// @Description and if fallback url is given, it will redirect to fallback url if it fails.
// @Description note that all urls must be url-encoded using QueryEscape in url package.
// @Param service query string false "subservice name"
// @Param sid query string false "subservice session id"
// @Param redirect query string false "redirect url"
// @Param fallback query string false "fallback url"
// @Success      200 "Ok, the client already has a valid session"
// @Success      201 "Created, the client has no valid session, so new session is created"
// @Success      303 "See Other, if redirect url is given, it redirects whether success or not"
// @Failure	  400 "Bad Request"
// @Failure 500 "Internal Server Error"
func (ac authApiController) HusSessionHandler(c echo.Context) error {
	service := c.QueryParam("service")
	sessionID := c.QueryParam("sid")
	redirectURL := c.QueryParam("redirect")
	fallbackURL := c.QueryParam("fallback")
	if fallbackURL == "" {
		fallbackURL = redirectURL
	}

	if (service != "" || sessionID != "" || redirectURL != "") && (service == "" || sessionID == "" || redirectURL == "") {
		return c.String(http.StatusBadRequest, "service, sid, redirect should be given all together or none")
	}

	var subservice common.ServiceDomain
	var ok bool

	// if service name is given, check if it is.
	// if not, return error.
	if service != "" {
		subservice, ok = common.Subservice[service]
		// if the service name is not registered, return error.
		if !ok {
			return c.String(http.StatusBadRequest, "no such service")
		}
	}

	// sessionID to UUID
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return c.String(http.StatusBadRequest, "sid is not valid")
	}

	// url decode
	redirectURL, err1 := url.QueryUnescape(redirectURL)
	fallbackURL, err2 := url.QueryUnescape(fallbackURL)
	if err1 != nil || err2 != nil {
		return c.String(http.StatusBadRequest, "url is not valid")
	}

	// when there's no valid Hus session, create new one depending on this flag.
	var createFlag bool

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

	// when there's no valid Hus session
	if createFlag {
		/* NEW HUS SESSION CREATION */
		CreateHusSessionParams := session.CreateHusSessionParams{
			Ctx:     c.Request().Context(),
			Dbc:     ac.dbClient,
			Service: &subservice,
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
