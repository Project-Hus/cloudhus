package auth

type SessionCheckBody struct {
	Sid string `json:"sid"`
	Uid string `json:"uid"`
}
