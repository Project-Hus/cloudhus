package auth

import (
	"hus-auth/ent"

	"github.com/labstack/echo/v4"
)

type AuthApiController struct {
	Client *ent.Client
}

func (c AuthApiController) AuthApiController(controller *AuthApiController) *echo.Echo {
	api := echo.New()

	api.POST("/auth/google", GoogleAuthHandler)

	return api
}
