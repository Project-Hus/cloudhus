package auth

import (
	"fmt"
	"hus-auth/ent"
	"hus-auth/helper"
	"hus-auth/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TokenEmbeddingHandler godoc
// @Router       /hus [post]
// @Summary      processes google auth and redirect with refresh token in url.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the refresh token will be expired in 7 days.
// @Tags         auth
// @Accept       json
// @Param        jwt body string true "Google ID token"
// @Success      301 "to /auth/{token_string}"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
func (ac authApiController) TokenEmbeddingHandler(c echo.Context) error {
	// get refresh token from header
	refreshToken := c.Request().Header.Get("Authorization")
	// validate refresh token
	_, err := service.ValidateRefreshToken(c.Request().Context(), ac.Client, refreshToken)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	// set cookie with refresh token
	cookie := &http.Cookie{
		Name:     "hus-refresh-token",
		Value:    refreshToken,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	fmt.Println("yo")
	// get whole set-cookie from echo context
	cookies := c.Cookies()

	for _, cookie := range cookies {
		fmt.Println(cookie.Name + cookie.Value)
	}

	return c.NoContent(http.StatusOK)
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

	err = ac.Client.HusSession.DeleteOneID(suuid).Exec(c.Request().Context())
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
