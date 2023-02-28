package api

import (
	"hus-auth/api/auth"

	"github.com/labstack/echo/v4"
)

func AuthApiController() *echo.Echo {
	api := echo.New()

	api.POST("/auth/google", auth.GoogleAuthHandler)

	return api
}
