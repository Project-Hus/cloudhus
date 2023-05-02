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
var Subservice = map[string]ServiceDomain{}

func init() {
	husenv, ok := os.LookupEnv("HUS_ENV")
	if !ok {
		log.Fatal("HUS_ENV is not set")
	}
	// at production level, we use actual domain names
	switch husenv {
	case "production":
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
	case "development":
		Subservice["cloudhus"] = ServiceDomain{
			Domain: Domain{
				Name: "cloudhus",
				URL:  "http://localhost",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "http://localhost",
				},
				"api": {
					Name: "api",
					URL:  "http://localhost",
				},
			},
		}
		Subservice["lifthus"] = ServiceDomain{
			Domain: Domain{
				Name: "lifthus",
				URL:  "http://localhost:3000",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "http://host.docker.internal:9100",
				},
				"api": {
					Name: "api",
					URL:  "http://host.docker.internal:9100",
				},
			},
		}
		Subservice["surfhus"] = ServiceDomain{
			Domain: Domain{
				Name: "surfhus",
				URL:  "https://localhost:3001",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "https://host.docker.internal:9200",
				},
				"api": {
					Name: "api",
					URL:  "https://host.docker.internal:9200",
				},
			},
		}
		// at native Go environment, we directly access localhost
	case "native":
		Subservice["cloudhus"] = ServiceDomain{
			Domain: Domain{
				Name: "cloudhus",
				URL:  "http://localhost:9001",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "http://localhost:9001",
				},
				"api": {
					Name: "api",
					URL:  "http://localhost:9002",
				},
			},
		}
		Subservice["lifthus"] = ServiceDomain{
			Domain: Domain{
				Name: "lifthus",
				URL:  "http://localhost:3000",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "http://localhost:9101",
				},
				"api": {
					Name: "api",
					URL:  "http://localhost:9102",
				},
			},
		}
		Subservice["surfhus"] = ServiceDomain{
			Domain: Domain{
				Name: "surfhus",
				URL:  "http://localhost:3001",
			},
			Subdomains: map[string]Domain{
				"auth": {
					Name: "auth",
					URL:  "http://localhost:9201",
				},
				"api": {
					Name: "api",
					URL:  "http://localhost:9202",
				},
			},
		}
	default:
		log.Fatal("HUS_ENV is not set properly. (production|development|native)")
	}
}
