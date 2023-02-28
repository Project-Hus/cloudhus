package auth

import (
	"context"
	"fmt"
	db "hus-auth/db/user"
	"hus-auth/ent"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"

	_ "net/http/httputil"
)

// GoogleAuthHandler godoc
// @Summary      Process google auth
// @Description  validates the google ID token and redirects with hus auth token query.
// @Tags         auth
// @Accept       json
// @Param        jwt path string true "Google ID token"
// @Success      301 "to /token"
// @Failure      401 "Unauthorized"
// @Router       /auth/google [post]
func (ac authApiController) GoogleAuthHandler(c echo.Context) error {
	credential := c.FormValue("credential")

	const clientID = "199526293983-r0b7tpmbpcc8nb786v261e451i2vihu3.apps.googleusercontent.com"

	payload, err := idtoken.Validate(context.TODO(), credential, clientID)
	if err != nil {
		log.Println("Invalid ID token: %w", err)
		return c.NoContent(http.StatusUnauthorized)
	}
	// Check that the user's ID token was intended for your application.
	if payload.Audience != clientID {
		log.Println("Invalid client ID")
		return c.NoContent(http.StatusUnauthorized)
	}

	sub := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)
	email_verified := payload.Claims["email_verified"].(bool)
	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)
	given_name := payload.Claims["given_name"].(string)
	family_name := payload.Claims["family_name"].(string)

	// query the user if it's registered.
	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)
	if err != nil {
		log.Println("query user failed")
		return c.NoContent(http.StatusInternalServerError)
	}
	// there is no Hus account with this Google account.
	if u == nil {
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
			log.Println("creating user failed %w", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		log.Printf("New user registerd(%s)", email)
	}

	// create auth token and send with url query.
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString("test")

	fmt.Println(tokenString, err)

	return c.Redirect(http.StatusMovedPermanently, "http://localhost:3000/")
}
