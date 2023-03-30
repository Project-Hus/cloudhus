package common

import (
	"os"
)

type ServiceDomain struct {
	Domain     Domain
	Subdomains map[string]Domain `json:"subdomains"`
}

type Domain struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// later it would be good to have this info in database and to initialize this map when the server starts
var Subservice = map[string]ServiceDomain{
	"cloudhus": {
		Domain: Domain{
			Name: "cloudhus",
			URL:  "https://cloudhus.com",
		},
		Subdomains: map[string]Domain{
			"auth": {
				Name: "auth",
				URL:  "https://auth.cloudhus.com",
			},
			"api": {
				Name: "api",
				URL:  "https://api.cloudhus.com",
			},
		},
	},
	"lifthus": {
		Domain: Domain{
			Name: "lifthus",
			URL:  "https://lifthus.com",
		},
		Subdomains: map[string]Domain{
			"auth": {
				Name: "auth",
				URL:  "https://auth.lifthus.com",
			},
			"api": {
				Name: "api",
				URL:  "https://api.lifthus.com",
			},
		},
	},
	"surfhus": {
		Domain: Domain{
			Name: "surfhus",
			URL:  "https://surfhus.com",
		},
		Subdomains: map[string]Domain{
			"auth": {
				Name: "auth",
				URL:  "https://auth.surfhus.com",
			},
			"api": {
				Name: "api",
				URL:  "https://api.surfhus.com",
			},
		},
	},
}

func init() {
	goenv := os.Getenv("GOENV")
	if cloudhus, ok := Subservice["cloudhus"]; ok {
		if subdomian, ok := cloudhus.Subdomains["auth"]; ok {
			subdomian.URL = "abc"
		}
	}
	if goenv == "development" {

	} else if goenv == "native" {

	}
}
