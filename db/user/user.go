package db

import (
	"context"
	"fmt"
	"hus-auth/ent"
	"hus-auth/ent/user"
	"log"
)

// CreateUserFromGoogle takes google ID token data and register new user to Project-Hus network.
func CreateUserFromGoogle(ctx context.Context, client *ent.Client, gu ent.User) (*ent.User, error) {
	u, err := client.User.
		Create().SetGoogleSub(gu.GoogleSub).SetEmail(gu.Email).SetEmailVerified(gu.EmailVerified).
		SetName(gu.Name).SetGoogleProfilePicture(gu.GoogleProfilePicture).SetFamilyName(gu.FamilyName).
		SetGivenName(gu.GivenName).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

// QuerUserByGoogleSub takes Google's sub and check if the user is registered.
func QueryUserByGoogleSub(ctx context.Context, client *ent.Client, sub string) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.GoogleSub(sub)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	return u, nil
}
