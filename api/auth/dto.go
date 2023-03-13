package auth

type HusSessionCheckBody struct {
	Sid string `json:"sid"`
	Uid string `json:"uid"`
}
