package auth

import (
	"hus-auth/common/db"
	"hus-auth/common/helper"
	"hus-auth/ent/hussession"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (ac authApiController) SSODemoHandler(c echo.Context) error {
	hst, err := c.Cookie("hus_st")
	if err != nil && err != http.ErrNoCookie {
		return c.HTML(http.StatusInternalServerError, basicHTMLGenerator("reading cookie failed T.T"))
	} else if err == http.ErrNoCookie {
		return c.HTML(http.StatusOK, basicHTMLGenerator("Hus session not established OoO"))
	}

	claims, expired, err := helper.ParseJWTWithHMAC(hst.Value)
	if err != nil {
		return c.HTML(http.StatusBadRequest, basicHTMLGenerator("Hus session is invalid -_-?"))
	} else if expired {
		return c.HTML(http.StatusBadRequest, basicHTMLGenerator("Hus session is expired -_o"))
	}

	hsid := claims["sid"].(string)
	hsuuid, err := uuid.Parse(hsid)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, basicHTMLGenerator("Hus session is invalid -_-?"))
	}

	hs, err := db.Client.HusSession.Query().Where(hussession.ID(hsuuid)).WithUser().Only(c.Request().Context())
	if err != nil {
		return c.HTML(http.StatusInternalServerError, basicHTMLGenerator("Hus session is invalid -_-?"))
	}

	u := hs.Edges.User

	if u == nil {
		return c.HTML(http.StatusOK, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>SSO Demonstration</title>
		</head>
		<body>
			<h1>SSO Demonstration</h1>
			<p>You are not signed in!</p>
			<p>Come and join us!</p>
			<p><a href="https://www.lifthus.com">Lifthus</a></p>
		</body>
		</html>
		`)
	}

	return c.HTML(http.StatusOK, `
	<!DOCTYPE html>
		<html>
		<head>
			<title>SSO Demonstration</title>
		</head>
		<body>
			<h1>SSO Demonstration</h1>
			<p>Hi!`+u.GivenName+`</p>
			<>Thank you for joining us!</p>
			<p><a href="https://www.lifthus.com">Lifthus</a></p>
		</body>
		</html>
	`)
}

func basicHTMLGenerator(message string) string {
	return `
	<!DOCTYPE html>
	<html>
	<head>
		<title>SSO Demonstration</title>
	</head>
	<body>
		<h1>SSO Demonstration</h1>
		<p>` + message + `</p>
	</body>
	</html>
	`
}
