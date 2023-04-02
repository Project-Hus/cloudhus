package hus

import (
	"hus-auth/ent"
	"net/http"
	"os"
)

var GoogleClientID = ""
var HusSecretKey = ""

var Host = ""
var URL = ""
var Origins = []string{}
var AuthCookieDomain = ""
var AuthURL = ""
var ApiURL = ""
var CookieSecure = false
var SameSiteMode = http.SameSiteNoneMode

var LifthusURL = "http://localhost:3000"

func InitHusVars(goenv string, _ *ent.Client) {
	//common
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	HusSecretKey = os.Getenv("HUS_SECRET_KEY")
	if goenv == "production" {
		Host = "cloudhus.com"
		URL = "https://cloudhus.com"
		Origins = []string{"https://cloudhus.com", "https://lifthus.com", "https://surfhus.com", "http://localhost:3000",
			"https://www.cloudhus.com", "https://www.lifthus.com", "https://www.surfhus.com"}
		AuthCookieDomain = "auth.cloudhus.com"
		AuthURL = "https://auth.cloudhus.com"
		ApiURL = "https://api.cloudhus.com"
		CookieSecure = true
		SameSiteMode = http.SameSiteNoneMode
	} else { // development or native
		Host = "localhost:9090"
		URL = "http://localhost:9090"
		Origins = []string{"http://localhost:3000", "http://localhost:9090", "http://localhost:9091"}
		AuthCookieDomain = ""
		AuthURL = "http://localhost:9090"
		ApiURL = "http://localhost:9090"
		CookieSecure = false
		SameSiteMode = http.SameSiteLaxMode
	}
	return
}
