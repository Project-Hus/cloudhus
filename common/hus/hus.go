package hus

import (
	"hus-auth/ent"
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

var LifthusURL = "http://localhost:3000"

func InitHusVars(_ *ent.Client) {
	goenv := os.Getenv("GOENV")
	//common
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if goenv == "production" {
		Host = "cloudhus.com"
		URL = "https://cloudhus.com"
		Origins = []string{"https://cloudhus.com", "https://lifthus.com", "https://surfhus.com"}
		AuthCookieDomain = "auth.cloudhus.com"
		AuthURL = "https://auth.cloudhus.com"
		ApiURL = "https://api.cloudhus.com"
		HusSecretKey = os.Getenv("HUS_SECRET_KEY")
	} else { // development
		Host = "localhost:9090"
		URL = "http://localhost:9090"
		Origins = []string{"http://localhost:3000", "http://localhost:9090", "http://localhost:9091"}
		AuthCookieDomain = ""
		AuthURL = "http://localhost:9090"
		ApiURL = "http://localhost:9090"
		HusSecretKey = os.Getenv("HUS_SECRET_KEY")
	}
}
