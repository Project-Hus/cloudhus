package auth

import (
	"bytes"
	"encoding/json"
	"hus-auth/common"
	"hus-auth/ent"
	"hus-auth/ent/hussession"
	"hus-auth/helper"
	"log"
	"net/http"
	"os"

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
func (ac authApiController) SessionCheckHandler(c echo.Context) error {
	// get service name and sid from path
	service := c.Param("service")
	lifthus_sid := c.Param("sid")

	subservice, ok := common.Subservice[service]
	// if the service name is not registered, return 404
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	// get hus_st from cookie
	hus_st, _ := c.Cookie("hus_st")
	// no valid st cookie, then return 401
	if hus_st == nil {
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
	scb := SessionCheckBody{lifthus_sid, hus_uid}
	scbBytes, _ := json.Marshal(scb)
	buff := bytes.NewBuffer(scbBytes)
	resp, err := http.Post(subservice.URL+"/session/check/", "application/json", buff)
	if err != nil {
		log.Println("[F] session transfering to "+subservice.Name+" failed: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return c.Redirect(http.StatusPermanentRedirect, subservice.URL+"/session/check/")
	} else {
		log.Println("[F] an error occured from " + subservice.Name)
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
