package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"

	_ "net/http/httputil"
)

// GoogleLoginHandler godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      301 "to /"
// @Failure      400
// @Router       /sign [post]
func GoogleLoginHandler(c echo.Context) error {
	credential := c.FormValue("credential")

	const clientID = "199526293983-r0b7tpmbpcc8nb786v261e451i2vihu3.apps.googleusercontent.com"

	payload, err := idtoken.Validate(context.TODO(), credential, clientID)
	if err != nil {
		// Handle any errors that occur while verifying the ID token.
		log.Fatalf("Invalid ID token: %v", err)
	}
	// Check that the user's ID token was intended for your application.
	if payload.Audience != clientID {
		log.Fatalf("Invalid client ID")
	}

	sub := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)
	email_verified := payload.Claims["email_verified"].(bool)
	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)
	given_name := payload.Claims["given_name"].(string)
	family_name := payload.Claims["family_name"].(string)

	fmt.Println(sub, email, email_verified, name, picture, given_name, family_name)

	return c.Redirect(http.StatusMovedPermanently, "http://localhost:3000/")
}
