package auth

import (
	"hus-auth/db"
	"hus-auth/helper"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

// GoogleAuthHandler godoc
// @Router       /auth/google [post]
// @Summary      processes google auth and redirect with refresh token in url.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the refresh token will be expired in 7 days.
// @Tags         auth
// @Accept       json
// @Param        jwt body string true "Google ID token"
// @Success      301 "to /auth/{token_string}"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
func (ac authApiController) GoogleAuthHandler(c echo.Context) error {
	// client ID that Google issued to lifthus
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	// credential sent from Google
	credential := c.FormValue("credential")
	// validating and parsing Google ID token
	payload, err := idtoken.Validate(c.Request().Context(), credential, clientID)
	if err != nil {
		log.Println("[F] Invalid ID token: %w", err)
		return c.NoContent(http.StatusUnauthorized)
	}
	// check if the user's ID token was intended for Lifthus
	if payload.Audience != clientID {
		log.Println("[F] Invalid client ID:", payload.Audience)
		return c.NoContent(http.StatusUnauthorized)
	}

	// Google's unique user ID
	sub := payload.Claims["sub"].(string)
	// check if the user is registered with Google
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		_, err := db.CreateUserFromGoogle(c.Request().Context(), ac.Client, *payload)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	// We checked or created if the Google user exists in Hus,
	// Now get user query again to create refresh token.
	u, err = db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// create and get refresh token
	refreshTokenSigned, err := db.CreateRefreshToken(c.Request().Context(), ac.Client, u.ID.String())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusMovedPermanently, os.Getenv("LIFTHUS_URL")+"/auth/"+refreshTokenSigned)
}

// AccessTokenRequestHandler godoc
// @Router       /auth/access [get]
// @Summary      gets refresh token in the header and returns access token in the cookie after validation.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the access token will be expired in 10 minutes.
// @Tags         auth
// @Param        jwt header string true "Refresh token"
// @Success      201 "Access token created"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
func (ac authApiController) AcessTokenRequestHandler(c echo.Context) error {
	// get refresh token from header
	refreshToken := c.Request().Header.Get("refresh_token")
	// validate refresh token
	refreshTokenValidated, err := helper.ValidateRefreshToken(c.Request().Context(), ac.Client, refreshToken)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// get user's uuid from refresh token
	uid := refreshTokenValidated["uid"].(string)

	// Create a new access token with 10 minutes expiration time.
	accessTokenSigned, err := helper.GetNewAccessToken(c.Request().Context(), ac.Client, uid)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = accessTokenSigned
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.NoContent(http.StatusCreated)
}
