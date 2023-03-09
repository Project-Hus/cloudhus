package auth

import (
	"fmt"
	"hus-auth/db"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

// GoogleAuthHandler godoc
// @Router       /social/google [post]
// @Summary      gets google IDtoken and redirect with refresh token in url.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the refresh token will be expired in 7 days.
// @Tags         auth
// @Accept       json
// @Param        jwt body string true "Google ID token"
// @Success      301 "to /auth/{token_string}"
// @Failure      301 "to /error"
func (ac authApiController) GoogleAuthHandler(c echo.Context) error {
	// client ID that Google issued to lifthus
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	// credential sent from Google
	credential := c.FormValue("credential")

	// get where the user redirected from
	fmt.Println(c.FormValue("redirect"))

	// validating and parsing Google ID token
	payload, err := idtoken.Validate(c.Request().Context(), credential, clientID)
	if err != nil {
		log.Println("[F] Invalid ID token: %w", err)
		return c.Redirect(http.StatusMovedPermanently, "http://"+c.Request().Host+"/error")
	}
	// check if the user's ID token was intended for Lifthus
	if payload.Audience != clientID {
		log.Println("[F] Invalid client ID:", payload.Audience)
		return c.Redirect(http.StatusMovedPermanently, "http://"+c.Request().Host+"/error")
	}

	// Google's unique user ID
	sub := payload.Claims["sub"].(string)
	// check if the user is registered with Google
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, "http://"+c.Request().Host+"/error")
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		_, err := db.CreateUserFromGoogle(c.Request().Context(), ac.Client, *payload)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "http://"+c.Request().Host+"/error")
		}
	}

	// We checked or created if the Google user exists in Hus,
	// Now get user query again to create refresh token.
	u, err = db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, "http://"+c.Request().Host+"/error")
	}

	// create and get refresh token
	refreshTokenSigned, err := db.CreateRefreshToken(c.Request().Context(), ac.Client, u.ID.String())
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, "http://"+c.Request().Host+"/error")
	}

	// set cookie for refresh token with 7 days expiration by struct literal
	cookie := &http.Cookie{
		Name:  "hus-refresh-token",
		Value: refreshTokenSigned,
		Path:  "/",
		//Secure:   true, // only sent over https
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		//Domain:   os.Getenv("COOKIE_DOMAIN"),
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	// redirects to lifthus.com/sign/token/{refresh_token}
	return c.Redirect(http.StatusMovedPermanently, os.Getenv("LIFTHUS_URL"))
}
