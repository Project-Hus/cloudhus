package auth

import (
	"hus-auth/ent"

	"github.com/labstack/echo/v4"
)

func NewAuthApiController(client *ent.Client) *echo.Echo {
	api := echo.New()

	authApiController := newAuthApiController(client)

	api.POST("/auth/google", authApiController.GoogleAuthHandler)

	return api
}

type authApis interface {
	GoogleAuthHandler(c echo.Context) error
}

type authApiController struct {
	Client *ent.Client
}

func newAuthApiController(client *ent.Client) authApis {
	return &authApiController{Client: client}
}
