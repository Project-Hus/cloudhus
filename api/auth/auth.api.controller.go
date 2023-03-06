package auth

import (
	"hus-auth/ent"

	"github.com/labstack/echo/v4"
)

// NewAuthApiController returns Echo comprising of auth api routes. instance to main.
func NewAuthApiController(client *ent.Client) *echo.Echo {
	authApi := echo.New()

	authApiController := newAuthApiController(client)

	authApi.POST("/auth/google", authApiController.GoogleAuthHandler)
	authApi.GET("/auth/cookie", authApiController.TokenEmbeddingHandler)
	authApi.GET("/auth/access", authApiController.AcessTokenRequestHandler)

	authApi.GET("/auth/refresh", authApiController.RefreshTokenRequestHandler)

	return authApi
}

// authApis interface defines what auth api has to handle
type authApis interface {
	GoogleAuthHandler(c echo.Context) error
	TokenEmbeddingHandler(c echo.Context) error
	AcessTokenRequestHandler(c echo.Context) error

	RefreshTokenRequestHandler(c echo.Context) error
}

// authApiController defines what auth api has to have and implements authApis interface at service file.
type authApiController struct {
	Client *ent.Client
}

// newAuthApiController returns a new authApiController that implements every auth api features.
func newAuthApiController(client *ent.Client) authApis {
	return &authApiController{Client: client}
}
