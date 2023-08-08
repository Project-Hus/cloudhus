package auth

import (
	"hus-auth/common/db"
	"hus-auth/common/helper"
	"hus-auth/ent/hussession"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SSODemoHandler godoc
// @Tags         auth
// @Router /demo/sso [get]
// @Summary shows the SSO feature between Cloudhus and Lifthus.
// @Success      200 "Ok, session is well-handled"
// @Failure      400 "Bad Request, something's wrong"
// @Failure      500 "Internal Server Error, something's wrong"
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
		return c.HTML(http.StatusOK, basicHTMLTemplate(`
			<div class="container">
			<h1>SSO Demonstration</h1>
			<img src="https://avatars.githubusercontent.com/u/124238598?s=200&v=4" alt="Project Hus Logo">
			<p>You are not signed in!</p>
			<p>Come and join us!</p>
			<h3><a href="https://www.lifthus.com" target="_blank">Lifthus</a></h3>
			</div>
		`))
	}

	var profImgURL string
	profImgURLp := u.ProfileImageURL
	if profImgURLp == nil {
		profImgURL = "https://avatars.githubusercontent.com/u/124238598?s=200&v=4"
	} else {
		profImgURL = *profImgURLp
	}

	return c.HTML(http.StatusOK, basicHTMLTemplate(`
			<div class="container">
			<h1>SSO Demonstration</h1>
			<img src="`+profImgURL+`" alt="Your profile image">
			<p>Hi! <b>`+u.GivenName+`</b>,</p>
			<p>and your family name is... <b>`+u.FamilyName+`</b>!</p>
			<p>Thank you for joining us!</p>
			<h3><a href="https://www.lifthus.com" target="_blank">Lifthus</a></h3>
			</div>
	`))
}

func basicHTMLGenerator(message string) string {
	return basicHTMLTemplate(`
	<div class="container">
	<h1>SSO Demonstration</h1>
	<img src="https://avatars.githubusercontent.com/u/124238598?s=200&v=4" alt="Project Hus Logo">
		<p>` + message + `</p>
		<h3><a href="https://www.lifthus.com" target="_blank">Lifthus</a></h3>
	</div>
	`)
}

func basicHTMLTemplate(content string) string {
	return `
	<!DOCTYPE html>
	<html>
	<head>
		<title>SSO Demonstration</title>
		` + commonStyle + `
	</head>
	<body>
	` + content + `
	</body>
	</html>
	`
}

var commonStyle = `
<style>
        img {
            width: 300px;
            height: 300px;
            object-fit: contain;
        }
		body {
			margin: 0;
			padding: 0;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			background-color: #363e50;
		  }
		.container {
			text-align: center;
			color: #9298e2;
		}
		a {
			color: #5dd1f1
		}
		a:visited {
			color: #03baec
		}
</style>
`
