package db

import (
	"context"
	"fmt"
	"hus-auth/dto"
	"hus-auth/ent"
	"hus-auth/ent/user"
	"log"
)

func CreateUser(ctx context.Context, client *ent.Client, gu dto.GoogleUser) (*ent.User, error) {
	u, err := client.User.
		Create().SetGoogleSub(gu.Sub).SetEmail(gu.Email).SetEmailVerified(gu.Email_verified).
		SetName(gu.Name).SetGoogleProfilePicture(gu.Picture).SetFamilyName(gu.Family_name).
		SetGivenName(gu.Given_name).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client, sub string) (*ent.User, error) {
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
