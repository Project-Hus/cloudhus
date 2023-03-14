package auth

import (
	"bytes"
	"encoding/json"
	"hus-auth/common"
	"hus-auth/db"
	"hus-auth/ent"
	"hus-auth/ent/hussession"
	"hus-auth/helper"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SessionCheckHandler godoc
// @Router /session/check/:service/:sid [post]
// @Summary accepts sid and service name to check if the session is valid.
// @Description  checks the hus session and tell subservice server if the session is valid.
// @Tags         auth
// @Param service path string true "subservice name"
// @Param sid path string true "session id"
// @Param        jwt header string false "Hus session token in cookie"
// @Success      200 "Ok"
// @Failure      401 "Unauthorized"
func (ac authApiController) HusSessionCheckHandler(c echo.Context) error {
	// get service name and sid from path
	service := c.Param("service")
	lifthus_sid := c.Param("sid")

	subservice, ok := common.Subservice[service]
	// if the service name is not registered, return 404
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	// get hus_st from cookie
	hus_st, err := c.Cookie("hus_st")
	// no valid st cookie, then return 401
	if hus_st == nil || err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// check if the session is valid
	claims, exp, err := helper.ParseJWTwithHMAC(hus_st.Value)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	} else if exp {
		// if the st is expired, then return 401.
		return c.NoContent(http.StatusUnauthorized)
	}
	// if the purpose is not hus_session, then return 401.
	if claims["purpose"].(string) != "hus_session" {
		return c.NoContent(http.StatusUnauthorized)
	}

	hus_sid := claims["sid"].(string)
	hus_uid := claims["uid"].(string)

	hus_sid_uuid, err := uuid.Parse(hus_sid)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// check if the hus session is not revoked querying the database with hus_sid.
	hs, err := ac.dbClient.HusSession.Query().Where(hussession.ID(hus_sid_uuid)).Only(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	} else if hs.UID.String() != hus_uid { // if the user ID is not matched, then return 401.
		return c.NoContent(http.StatusUnauthorized)
	}

	// now we know that the hus session is valid, so we tell the subservice server that the session is valid with uid.
	// make http request to subservice server
	u, err := db.QueryUserByUID(c.Request().Context(), ac.dbClient, hus_uid)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	bd := ""
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
		log.Println("[F] marshaling body to lifthus failed: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	buff := bytes.NewBuffer(scbBytes)

	// with ac.httpClient, transfer the validation result to subservice auth server.
	req, err := http.NewRequest("POST", subservice.Subdomains["auth"].URL+"/hus/session/check", buff)
	if err != nil {
		log.Println("[F] session transfering to "+subservice.Domain.Name+" failed: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")
	// send the request
	resp, err := ac.httpClient.Do(req)
	if err != nil {
		log.Println("[F] session transfering to "+subservice.Domain.Name+" failed: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return c.NoContent(http.StatusOK)
		//return c.Redirect(http.StatusPermanentRedirect, subservice.Subdomains["auth"].URL+"/session/check")
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
