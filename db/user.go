package db

import (
	"context"
	"hus-auth/ent"
	"hus-auth/ent/user"
	"log"

	"google.golang.org/api/idtoken"
)

// CreateUserFromGoogle takes google ID token data and register new user to Project-Hus network.
func CreateUserFromGoogle(ctx context.Context, client *ent.Client, payload idtoken.Payload) (*ent.User, error) {
	// Google user information to use as Hus user information
	sub := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)
	emailVerified := payload.Claims["email_verified"].(bool)
	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)
	givenName := payload.Claims["given_name"].(string)
	familyName := payload.Claims["family_name"].(string)
	u, err := client.User.
		Create().SetGoogleSub(sub).SetEmail(email).SetEmailVerified(emailVerified).
		SetName(name).SetProfilePictureURL(picture).SetFamilyName(familyName).
		SetGivenName(givenName).Save(ctx)
	if err != nil {
		log.Printf("[F] creating user failed:%v", err)
		return nil, err
	}
	log.Printf("new user(%s) is registered by Google", u.ID)
	return u, nil
}

// QuerUserByGoogleSub takes Google's sub and check if the user is registered.
func QueryUserByGoogleSub(ctx context.Context, client *ent.Client, sub string) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.GoogleSub(sub)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		log.Printf("[F] querying user failed:%v", err)
		return nil, err
	}
	return u, err
}
