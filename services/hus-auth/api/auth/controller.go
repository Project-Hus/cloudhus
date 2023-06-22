package auth

import (
	"hus-auth/ent"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthApiControllerParams struct {
	DbClient   *ent.Client
	HttpClient *http.Client
}

// NewAuthApiController returns Echo comprising of auth api routes. instance to main.
func NewAuthApiController(authApi *echo.Echo, params AuthApiControllerParams) *echo.Echo {
	authApiController := newAuthApiController(params)

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

	// social login services
	authApi.POST("/auth/social/google/:service", authApiController.GoogleAuthHandler)

	// session services
	authApi.POST("/auth/session/check/:service/:sid", authApiController.HusSessionCheckHandler)
	authApi.DELETE("/auth/session/revoke", authApiController.SessionRevocationHandler)

	// social login services

	// Hus session services
	authApi.GET("/auth/hussession", authApiController.HusSessionHandler)
	authApi.GET("/auth/hussession/:token", authApiController.SessionConnectionHandler)

	return authApi
}

// newAuthApiController returns a new authApiController that implements every auth api features.
func newAuthApiController(params AuthApiControllerParams) authApis {
	return &authApiController{dbClient: params.DbClient, httpClient: params.HttpClient}
}

// authApiController defines what auth api has to have and implements authApis interface at service file.
type authApiController struct {
	dbClient   *ent.Client
	httpClient *http.Client
}

// authApis interface defines what auth api has to handle
type authApis interface {

	/* client side api */
	GoogleAuthHandler(c echo.Context) error

	HusSessionCheckHandler(c echo.Context) error
	SessionRevocationHandler(c echo.Context) error

	// Hus session services
	HusSessionHandler(c echo.Context) error
	SessionConnectionHandler(c echo.Context) error
}
