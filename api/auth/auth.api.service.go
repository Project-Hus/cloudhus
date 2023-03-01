package auth

import (
	"context"
	"fmt"
	"hus-auth/db"
	"hus-auth/ent"
	"hus-auth/helper"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"

	_ "net/http/httputil"
)

// GoogleAuthHandler godoc
// @Router       /auth/google [post]
// @Summary      processes google auth and redirect with refresh token in url.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
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
	payload, err := idtoken.Validate(context.TODO(), credential, clientID)
	if err != nil {
		log.Println("Invalid ID token: %w", err)
		return c.NoContent(http.StatusUnauthorized)
	}
	// check if the user's ID token was intended for Lifthus
	if payload.Audience != clientID {
		log.Println("Invalid client ID")
		return c.NoContent(http.StatusUnauthorized)
	}

	// Google's unique user ID
	sub := payload.Claims["sub"].(string)
	// check if the user is registered with Google
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		log.Println("checking Google user failed:%w", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	// create one if there is no Hus account with this Google account
	if u == nil {
		// Google user information to use as Hus user information
		email := payload.Claims["email"].(string)
		email_verified := payload.Claims["email_verified"].(bool)
		name := payload.Claims["name"].(string)
		picture := payload.Claims["picture"].(string)
		given_name := payload.Claims["given_name"].(string)
		family_name := payload.Claims["family_name"].(string)
		new_user := ent.User{
			GoogleSub:            sub,
			Email:                email,
			EmailVerified:        email_verified,
			Name:                 name,
			GoogleProfilePicture: picture,
			GivenName:            given_name,
			FamilyName:           family_name,
		}
		_, err := db.CreateUserFromGoogle(c.Request().Context(), ac.Client, new_user)
		if err != nil {
			log.Println("creating Google user failed:%w", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		log.Printf("New google user(%s) registerd", sub)
	}

	// We checked or created if the Google user exists in Hus,
	// Now get user query again to create refresh token.
	u, err = db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		log.Println("query user failed")
		return c.NoContent(http.StatusInternalServerError)
	}

	// create and get refresh token
	refreshTokenSigned, err := db.CreateRefreshToken(c.Request().Context(), ac.Client, u.ID.String())
	if err != nil {
		log.Println("creating signed refresh token failed:%w", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusMovedPermanently, os.Getenv("LIFTHUS_URL")+"/auth/"+refreshTokenSigned)
}

// GoogleAuthHandler godoc
// @Router       /auth/access [get]
// @Summary      gets refresh token in the header and returns access token after validation.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
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
		log.Println("validating refresh token failed:%w", err)
		return c.NoContent(http.StatusUnauthorized)
	}

	// get user's uuid from refresh token
	uid := refreshTokenValidated["uid"].(string)

	// Create a new access token with 10 minutes expiration time.
	aid := uuid.New().String()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aid":     aid,                       // refresh token's uuid
		"purpose": "access",                  // purpose
		"iss":     "https://api.lifthus.com", // issuer
		"uid":     uid,                       // user's uuid
		"iat":     time.Now().Unix(),         // issued at
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := accessToken.SignedString(os.Getenv("HUS_AUTH_TOKEN_KEY"))

	fmt.Println(tokenString, err)

}
