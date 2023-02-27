package main

import (
	"context"
	"net/http"
	"os"

	"fmt"
	"log"

	"hus-auth/ent"
	"hus-auth/ent/user"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"google.golang.org/api/idtoken"

	_ "github.com/go-sql-driver/mysql"
)

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

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
