package auth

import (
	"hus-auth/ent"

	"github.com/labstack/echo/v4"
)

// authApis interface defines what auth api has to handle
type authApis interface {
	// get google ID token and set hus session cookie
	GoogleAuthHandler(c echo.Context) error
	TokenEmbeddingHandler(c echo.Context) error
	AcessTokenRequestHandler(c echo.Context) error

	RefreshTokenRequestHandler(c echo.Context) error
}

// authApiController defines what auth api has to have and implements authApis interface at service file.
type authApiController struct {
	Client *ent.Client
}

// NewAuthApiController returns Echo comprising of auth api routes. instance to main.
func NewAuthApiController(client *ent.Client) *echo.Echo {
	authApi := echo.New()

	authApiController := newAuthApiController(client)

	authApi.POST("/google", authApiController.GoogleAuthHandler)
	authApi.GET("/hus", authApiController.TokenEmbeddingHandler)
	authApi.GET("/access", authApiController.AcessTokenRequestHandler)

	authApi.GET("/refresh", authApiController.RefreshTokenRequestHandler)

	return authApi
}

// newAuthApiController returns a new authApiController that implements every auth api features.
func newAuthApiController(client *ent.Client) authApis {
	return &authApiController{Client: client}
}
