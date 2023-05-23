package auth

import (
	"strconv"
	"strings"

	"fmt"
	"hus-auth/common"
	"hus-auth/common/hus"

	"hus-auth/service/session"
	"log"
	"net/http"
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
		"sid":            lifthus_sid,
		"uid":            strconv.FormatUint(u.ID, 10),
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"name":           u.Name,
		"given_name":     u.GivenName,
		"family_name":    u.FamilyName,
		"birthdate":      bd,
		"exp":            time.Now().Add(time.Second * 10).Unix(),
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
