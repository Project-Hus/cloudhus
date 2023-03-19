package auth

import (
	"strings"

	"fmt"
	"hus-auth/common"
	"hus-auth/db"
	"hus-auth/ent"
	"hus-auth/helper"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		return c.String(http.StatusUnauthorized, "[F]session expired")
	}
	// if the purpose is not hus_session, then return 401.
	if claims["purpose"].(string) != "hus_session" {
		return c.String(http.StatusUnauthorized, "[F]wrong purpose")
	}

	hus_sid := claims["sid"].(string)
	hus_uid := claims["uid"].(string)
	hus_iat := claims["iat"].(float64)

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

	hscbJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":            lifthus_sid,
		"uid":            hus_uid,
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"name":           u.Name,
		"given_name":     u.GivenName,
		"family_name":    u.FamilyName,
		"birthdate":      bd,
		"exp":            time.Now().Add(time.Second * 10).Unix(),
	})

	hscbSigned, err := hscbJWT.SignedString([]byte(os.Getenv("HUS_SECRET_KEY")))
	if err != nil {
		err = fmt.Errorf("[F]signing jwt for %s failed: %w", service, err)
		log.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// with ac.httpClient, transfer the validation result to subservice auth server.
	req, err := http.NewRequest("POST", subservice.Subdomains["auth"].URL+"/hus/session/sign", strings.NewReader(hscbSigned))
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
		// if the session check is successful, renew the hus_st.
		nhst := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sid":     hus_sid,                              // session token's uuid
			"purpose": "hus_session",                        // purpose"
			"iss":     os.Getenv("HUS_AUTH_URL"),            // issuer
			"uid":     hus_uid,                              // user's uuid
			"iat":     hus_iat,                              // issued at
			"exp":     time.Now().Add(time.Hour * 1).Unix(), // expiration : an hour
		})
		nhstSigned, err := nhst.SignedString([]byte(os.Getenv("HUS_SECRET_KEY")))
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		nhstCookie := &http.Cookie{
			Name:     "hus_st",
			Value:    nhstSigned,
			Path:     "/",
			Domain:   os.Getenv("HUS_DOMAIN"),
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Secure:   false,
		}
		c.SetCookie(nhstCookie)
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
	if hus_st == nil || hus_st.Value == "" {
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

// function that unmarsahls the jwt token and returns the claims.
func ParseJWTwithHMAC(token string) (jwt.MapClaims, *jwt.Token, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("HUS_SECRET_KEY")), nil
	})
	return claims, parsedToken, err
}
