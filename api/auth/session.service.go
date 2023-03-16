package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hus-auth/common"
	"hus-auth/db"
	"hus-auth/ent"
	"hus-auth/helper"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// HusSessionCheckHandler godoc
// @Router /session/check/:service/:sid [post]
// @Summary accepts sid and service name to check if the session is valid.
// @Description  checks the hus session in cookie and tells the subservice server if the session is valid.
// @Tags         auth
// @Param service path string true "subservice name"
// @Param sid path string true "session id"
// @Success      200 "Ok"
// @Failure      401 "Unauthorized"
func (ac authApiController) HusSessionCheckHandler(c echo.Context) error {
	// get service name and sid from path
	service := c.Param("service")
	lifthus_sid := c.Param("sid")

	subservice, ok := common.Subservice[service]
	// if the service name is not registered, return 404
	if !ok {
		return c.String(http.StatusNotFound, "[F]no such service")
	}

	// get hus_st from cookie
	hus_st, err := c.Cookie("hus_st")
	// no valid st cookie, then return 401
	if err != nil || hus_st.Value == "" {
		return c.String(http.StatusUnauthorized, "[F]not sigend in")
	}

	// check if the session is valid
	claims, exp, err := helper.ParseJWTwithHMAC(hus_st.Value)
	if err != nil {
		log.Printf("%v(from /session/check/%s/:sid)", err, service)
		return c.String(http.StatusUnauthorized, "[F]invalid session")
	} else if exp {
		// if the st is expired, then return 401.
		return c.String(http.StatusUnauthorized, "[F]session expired")
	}
	// if the purpose is not hus_session, then return 401.
	if claims["purpose"].(string) != "hus_session" {
		return c.String(http.StatusUnauthorized, "[F]wrong purpose")
	}

	hus_sid := claims["sid"].(string)
	hus_uid := claims["uid"].(string)

	// check if the hus session is not revoked querying the database with hus_sid.
	_, err = db.QuerySessionBySID(c.Request().Context(), ac.dbClient, hus_sid)
	if err != nil {
		return c.String(http.StatusUnauthorized, "[F]invalid session")
	}

	// now we know that the hus session is valid, so we tell the subservice server that the session is valid with uid.
	u, err := db.QueryUserByUID(c.Request().Context(), ac.dbClient, hus_uid)
	if err != nil {
		return c.String(http.StatusUnauthorized, "[F]no such user")
	}

	var bd string
	if u.Birthdate != nil {
		bd = u.Birthdate.Format(time.RFC3339)
	}
	scb := HusSessionCheckBody{
		lifthus_sid,
		hus_uid,
		u.Email,
		u.EmailVerified,
		u.Name,
		u.GivenName,
		u.FamilyName,
		bd,
	}

	scbBytes, err := json.Marshal(scb)
	if err != nil {
		err = fmt.Errorf("[F]marshalling body for %s failed: %w", service, err)
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	buff := bytes.NewBuffer(scbBytes)

	// with ac.httpClient, transfer the validation result to subservice auth server.
	req, err := http.NewRequest("POST", subservice.Subdomains["auth"].URL+"/hus/session/sign", buff)
	if err != nil {
		err = fmt.Errorf("[F]session injection to "+subservice.Domain.Name+" failed:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// send the request
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("[F]session injection to "+subservice.Domain.Name+" failed:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return c.String(http.StatusOK, "session injection to "+subservice.Domain.Name+" success")
	} else {
		log.Println("[F] an error occured from " + subservice.Domain.Name + ":" + resp.Status)
		return c.NoContent(http.StatusInternalServerError)
	}
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

	err = ac.dbClient.HusSession.DeleteOneID(suuid).Exec(c.Request().Context())
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
