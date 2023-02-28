package auth

import (
	"context"
	"fmt"
	db "hus-auth/db/user"
	"net/http"

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
func (ac AuthApiController) GoogleAuthHandler(c echo.Context) error {
	credential := c.FormValue("credential")

	const clientID = "199526293983-r0b7tpmbpcc8nb786v261e451i2vihu3.apps.googleusercontent.com"

	payload, err := idtoken.Validate(context.TODO(), credential, clientID)
	if err != nil {
		// Handle any errors that occur while verifying the ID token.
		return c.String(401, fmt.Sprintf("Invalid ID token: %v", err))
	}
	// Check that the user's ID token was intended for your application.
	if payload.Audience != clientID {
		return c.String(401, "Invalid client ID")
	}

	sub := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)
	email_verified := payload.Claims["email_verified"].(bool)
	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)
	given_name := payload.Claims["given_name"].(string)
	family_name := payload.Claims["family_name"].(string)

	u, err := db.QueryUserByGoogleSub(c.Request().Context(), ac.Client, sub)

	fmt.Println(sub, email, email_verified, name, picture, given_name, family_name)

	return c.Redirect(http.StatusMovedPermanently, "http://localhost:3000/")
}
