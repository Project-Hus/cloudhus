package auth

import (
	"hus-auth/ent"
	"net/http"

	"github.com/labstack/echo/v4"
)

// authApis interface defines what auth api has to handle
type authApis interface {

	/* client side api */
	GoogleAuthHandler(c echo.Context) error

	HusSessionCheckHandler(c echo.Context) error
	SessionRevocationHandler(c echo.Context) error
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
func NewAuthApiController(authApi *echo.Echo, params AuthApiControllerParams) *echo.Echo {

	authApiController := newAuthApiController(params)

	authApi.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Cloudhus")
	})

	// social login services
	authApi.POST("/auth/social/google/:service", authApiController.GoogleAuthHandler)

	// session services
	authApi.POST("/auth/session/check/:service/:sid", authApiController.HusSessionCheckHandler)
	authApi.DELETE("/auth/session/revoke", authApiController.SessionRevocationHandler)

	return authApi
}

// newAuthApiController returns a new authApiController that implements every auth api features.
func newAuthApiController(params AuthApiControllerParams) authApis {
	return &authApiController{dbClient: params.DbClient, httpClient: params.HttpClient}
}
