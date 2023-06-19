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
	dbHost, ok1 := os.LookupEnv("HUS_AUTH_DB_HOST")
	dbPort, ok2 := os.LookupEnv("HUS_AUTH_DB_PORT")
	dbUser, ok3 := os.LookupEnv("HUS_AUTH_DB_USER")
	dbPassword, ok4 := os.LookupEnv("HUS_AUTH_DB_PASSWORD")
	dbName, ok5 := os.LookupEnv("HUS_AUTH_DB_NAME")
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		log.Fatal("DB credentials are not set")
	}
	// DB connection
	connectionPhrase := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	client, err := ent.Open("mysql", connectionPhrase)
	if err != nil {
		log.Print(" opening connection to mysql failed: %w", err)
		return nil, err
	}

	// Running the auto migration tool.
	if err := client.Schema.Create(
		context.Background(),
	); err != nil {
		log.Print(" creating schema resources failed: %w", err)
		return nil, err
	}

	return client, nil
}

// Rollback is a helper function to rollback transaction.
func Rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
