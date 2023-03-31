package common

import (
	"log"
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

// Subservice defines the name and url of each service.
// in production, we use actual domain names.
// in development, we use localhost and docker internal network with SAM.
// in native Go, we use localhost.
var Subservice map[string]ServiceDomain

func init() {
	goenv := os.Getenv("GOENV")
	// at production level, we use actual domain names
	if goenv == "production" {
		Subservice["cloudhus"] = ServiceDomain{
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
		}
		Subservice["lifthus"] = ServiceDomain{
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
		}
		Subservice["surfhus"] = ServiceDomain{
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
		}
		// at sam local, we have to access localhost with docker network
	} else if goenv == "development" {
		Subservice["cloudhus"] = ServiceDomain{
			Domain: Domain{
				Name: "cloudhus",
				URL:  "https://localhost:9090",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "https://localhost:9090",
				},
				"api": {
					Name: "api",
					URL:  "https://localhost:9090",
				},
			},
		}
		Subservice["lifthus"] = ServiceDomain{
			Domain: Domain{
				Name: "lifthus",
				URL:  "https://host.docker.internal:9091",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "https://host.docker.internal:9091",
				},
				"api": {
					Name: "api",
					URL:  "https://host.docker.internal:9091",
				},
			},
		}
		Subservice["surfhus"] = ServiceDomain{
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
		}
		// at native Go environment, we directly access localhost
	} else if goenv == "native" {
		Subservice["cloudhus"] = ServiceDomain{
			Domain: Domain{
				Name: "cloudhus",
				URL:  "https://localhost:9090",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "https://localhost:9090",
				},
				"api": {
					Name: "api",
					URL:  "https://localhost:9090",
				},
			},
		}
		Subservice["lifthus"] = ServiceDomain{
			Domain: Domain{
				Name: "lifthus",
				URL:  "https://localhost:9091",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "https://localhost:9091",
				},
				"api": {
					Name: "api",
					URL:  "https://localhost:9091",
				},
			},
		}
		Subservice["surfhus"] = ServiceDomain{
			Domain: Domain{
				Name: "surfhus",
				URL:  "https://localhost:9092",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "https://localhost:9092",
				},
				"api": {
					Name: "api",
					URL:  "https://localhost:9092",
				},
			},
		}
	} else {
		log.Fatal("GOENV is not set properly. production|development|native ")
	}
}
