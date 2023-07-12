package dto

type HusConnDto struct {
	Hsid string       `json:"hsid"`
	User *HusConnUser `json:"user"`
}

type HusConnUser struct {
	Uid             uint64  `json:"uid"`
	ProfileImageURL *string `json:"profile_image_url"`
	Email           string  `json:"email"`
	EmailVerified   bool    `json:"email_verified"`
	Name            string  `json:"name"`
	GivenName       string  `json:"given_name"`
	FamilyName      string  `json:"family_name"`
}
