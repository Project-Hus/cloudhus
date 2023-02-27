package main

import (
	"context"
	"net/http"
	"os"

	"fmt"
	"log"

	"hus-auth/ent"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"

	_ "github.com/go-sql-driver/mysql"

	_ "hus-auth/docs" // docs is generated by Swag CLI, you have to import it.

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Project-Hus auth server
// @version 0.0.0
// @description This is Project-Hus's root authentication server containing each user's UUID, which is unique for all hus services.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url lifthus531@gmail.com
// @contact.email lifthus531@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host lifthus.com
// @BasePath /auth
func main() {
	// set .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error lading .env file: %s", err)
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// DB connection
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUser, dbPassword, dbHost, dbPort, dbName)
	client, err := ent.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// Running the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/sign", func(c echo.Context) error {
		credential := c.FormValue("credential")

		const clientID = "199526293983-r0b7tpmbpcc8nb786v261e451i2vihu3.apps.googleusercontent.com"

		payload, err := idtoken.Validate(context.TODO(), credential, clientID)
		if err != nil {
			// Handle any errors that occur while verifying the ID token.
			log.Fatalf("Invalid ID token: %v", err)
		}
		// Check that the user's ID token was intended for your application.
		if payload.Audience != clientID {
			log.Fatalf("Invalid client ID")
		}

		sub := payload.Claims["sub"].(string)
		email := payload.Claims["email"].(string)
		email_verified := payload.Claims["email_verified"].(bool)
		name := payload.Claims["name"].(string)
		picture := payload.Claims["picture"].(string)
		given_name := payload.Claims["given_name"].(string)
		family_name := payload.Claims["family_name"].(string)

		fmt.Println(sub, email, email_verified, name, picture, given_name, family_name)

		return c.Redirect(http.StatusMovedPermanently, "http://localhost:3000/")
	})
	e.Logger.Fatal(e.Start(":9090"))
}
