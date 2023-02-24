package main

import (
	"log"
	"net/http"
	"os"

	"hus-auth/ent"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error lading .env file: %s", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	client, err := ent.Open("mysql", "lifthus_auth")

	e := echo.New()
	e.POST("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Lifthus!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
