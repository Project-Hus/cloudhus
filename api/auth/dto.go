package auth

type HusSessionCheckBody struct {
	Sid           string `json:"sid"`
	Uid           string `json:"uid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Birthdate     string `json:"birthdate"`
}
