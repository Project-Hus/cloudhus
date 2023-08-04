package hus

import (
	"hus-auth/ent"
	"log"
	"net/http"
	"os"
	"time"
)

var GoogleClientID = ""
var HusSecretKey = ""
var HusSecretKeyBytes []byte

var Host = ""
var URL = ""
var Origins = []string{}
var AuthCookieDomain = ""
var AuthURL = ""
var ApiURL = ""

var CookieSecure = true
var CookieSameSite = http.SameSiteLaxMode

var LifthusURL = "http://localhost:3000"

var Http *http.Client

func InitHusVars(husenv string, _ *ent.Client) {
	Http = &http.Client{Timeout: 5 * time.Second}

	ok1, ok2 := false, false
	GoogleClientID, ok1 = os.LookupEnv("GOOGLE_CLIENT_ID")
	HusSecretKey, ok2 = os.LookupEnv("HUS_SECRET_KEY")
	HusSecretKeyBytes = []byte(HusSecretKey)
	if !ok1 || !ok2 {
		log.Fatal("GOOGLE_CLIENT_ID or HUS_SECRET_KEY is not set")
	}
	switch husenv {
	case "production":
		Host = "cloudhus.com"
		URL = "https://cloudhus.com"
		Origins = []string{"https://cloudhus.com", "https://lifthus.com", "https://surfhus.com", "http://localhost:3000",
			"https://www.cloudhus.com", "https://www.lifthus.com", "https://www.surfhus.com"}
		AuthCookieDomain = "auth.cloudhus.com"
		AuthURL = "https://auth.cloudhus.com"
		ApiURL = "https://api.cloudhus.com"
	case "development":
		Host = "localhost:9000"
		URL = "http://localhost:9000"
		Origins = []string{"http://localhost:3000", "http://localhost:9100", "http://localhost:9200"}
		AuthCookieDomain = ""
		AuthURL = "http://localhost:9000"
		ApiURL = "http://localhost:9000"
		CookieSecure = false
	case "native":
		Host = "localhost:9001"
		URL = "http://localhost:9001"
		Origins = []string{
			"http://localhost:3000",
			"http://localhost:9001",
			"http://localhost:9002",
			"http://localhost:9101",
			"http://localhost:9102",
		}
		AuthCookieDomain = ""
		AuthURL = "http://localhost:9001"
		ApiURL = "http://localhost:9002"
	default:
		log.Fatal("HUS_ENV is not set properly. (production|development|native)")
	}
}
