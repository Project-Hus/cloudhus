package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"hus-auth/ent"
)

// ConncectToHusAuth returns hus_auth_db's ent client.
// you've got to close it with Close() in defer out of this function.
func ConnectToHusAuth() (*ent.Client, error) {
	dbHost := os.Getenv("HUS_DB_HOST")
	dbPort := os.Getenv("HUS_DB_PORT")
	dbUser := os.Getenv("HUS_AUTH_DB_USER")
	dbPassword := os.Getenv("HUS_AUTH_DB_PASSWORD")
	dbName := os.Getenv("HUS_AUTH_DB_NAME")

	// DB connection
	connectionPhrase := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	client, err := ent.Open("mysql", connectionPhrase)
	if err != nil {
		log.Print(" opening connection to mysql failed: %w", err)
		return nil, err
	}

	// Running the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Print(" creating schema resources failed: %w", err)
		return nil, err
	}

	return client, nil
}
