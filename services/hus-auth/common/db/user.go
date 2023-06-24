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
		Create().
		SetGoogleSub(sub).
		SetEmail(email).
		SetEmailVerified(emailVerified).
		SetProfilePictureURL(picture).
		SetName(name).
		SetGivenName(givenName).
		SetFamilyName(familyName).
		SetProvider("google").
		Save(ctx)
	if err != nil {
		log.Printf("creating google user failed:%v", err)
		return nil, err
	}
	log.Printf("new user(%d) is registered by Google", u.ID)
	return u, nil
}

// QuerUserByGoogleSub takes Google's sub and check if the user is registered.
// it returns nil, nil if the user is not registered.
func QueryUserByGoogleSub(ctx context.Context, client *ent.Client, sub string) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.GoogleSub(sub)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		log.Printf("querying google user failed:%v", err)
		return nil, err
	}
	return u, err
}

// QueryUserByUID takes user's UID and returns user entity.
func QueryUserByUID(ctx context.Context, client *ent.Client, uid uint64) (*ent.User, error) {
	u, err := client.User.Query().Where(user.ID(uid)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		log.Printf("querying user failed:%v", err)
		return nil, err
	}
	return u, nil
}

// UUID version
// func QueryUserByUID(ctx context.Context, client *ent.Client, uid string) (*ent.User, error) {
// 	uid_uuid, err := uuid.Parse(uid)
// 	if err != nil {
// 		log.Println("[F] parsing uid failed:", err)
// 		return nil, err
// 	}
// 	u, err := client.User.Query().Where(user.ID(uid_uuid)).Only(ctx)
// 	if err != nil && !ent.IsNotFound(err) {
// 		log.Printf("[F] querying user failed:%v", err)
// 		return nil, err
// 	}
// 	return u, nil
// }
