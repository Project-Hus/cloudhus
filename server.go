package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Lifthus!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
