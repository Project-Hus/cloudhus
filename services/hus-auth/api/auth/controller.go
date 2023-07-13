package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewAuthApiController returns Echo comprising of auth api routes. instance to main.
func NewAuthApiController(authApi *echo.Echo) *echo.Echo {
	authApiController := newAuthApiController()

	authApi.GET("/auth", func(c echo.Context) error {
		return c.String(http.StatusOK, "\"모든 인증은 Cloudhus로 통한다\" -Cloudhus-")
	})

	authApi.GET("/auth/test/cookie", func(c echo.Context) error {
		// return all cookies as a string
		cookies := c.Cookies()
		cookiesString := "COOKIES:"
		for _, cookie := range cookies {
			cookiesString += cookie.Name + ": " + cookie.Value + "\n"
		}
		return c.String(http.StatusOK, cookiesString)
	})

	// social sign api
	authApi.POST("/auth/hus/sign/social/google", authApiController.GoogleAuthHandler) // from Client

	// Hus session api
	authApi.GET("/auth/hus", authApiController.HusSessionHandler)                       // from Client
	authApi.GET("/auth/hus/connect/:token", authApiController.SessionConnectionHandler) // from Subservice

	// sign out services
	authApi.PATCH("/auth/hus/signout", authApiController.SignOutHandler) // from Subservice

	return authApi
}

// newAuthApiController returns a new authApiController that implements every auth api features.
func newAuthApiController() authApis {
	return &authApiController{}
}

// authApiController defines what auth api has to have and implements authApis interface at service file.
type authApiController struct {
}

// authApis interface defines what auth api has to handle
type authApis interface {

	// social sign api
	GoogleAuthHandler(c echo.Context) error

	// Hus session api
	HusSessionHandler(c echo.Context) error
	SessionConnectionHandler(c echo.Context) error
	SignOutHandler(c echo.Context) error
}
