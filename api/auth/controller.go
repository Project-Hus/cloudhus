package auth

import (
	"hus-auth/ent"
	"net/http"

	"github.com/labstack/echo/v4"
)

// authApis interface defines what auth api has to handle
type authApis interface {
	// get google ID token and set hus session cookie
	GoogleAuthHandler(c echo.Context) error
	SessionRevocationHandler(c echo.Context) error

	TokenEmbeddingHandler(c echo.Context) error
	AcessTokenRequestHandler(c echo.Context) error

	RefreshTokenRequestHandler(c echo.Context) error
}

// authApiController defines what auth api has to have and implements authApis interface at service file.
type authApiController struct {
	dbClient   *ent.Client
	httpClient *http.Client
}

type AuthApiControllerParams struct {
	DbClient   *ent.Client
	HttpClient *http.Client
}

// NewAuthApiController returns Echo comprising of auth api routes. instance to main.
func NewAuthApiController(params AuthApiControllerParams) *echo.Echo {
	authApi := echo.New()

	authApiController := newAuthApiController(params)

	authApi.POST("/social/google/:service", authApiController.GoogleAuthHandler)
	authApi.DELETE("/session/revoke", authApiController.SessionRevocationHandler)

	authApi.GET("/hus", authApiController.TokenEmbeddingHandler)
	authApi.GET("/access", authApiController.AcessTokenRequestHandler)

	authApi.GET("/refresh", authApiController.RefreshTokenRequestHandler)

	return authApi
}

// newAuthApiController returns a new authApiController that implements every auth api features.
func newAuthApiController(params AuthApiControllerParams) authApis {
	return &authApiController{dbClient: params.DbClient, httpClient: params.HttpClient}
}
