package common

type subservice struct {
	Domain     Domain
	Subdomains map[string]Domain `json:"subdomains"`
}

type Domain struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// later it would be good to have this info in a database and to initialize this map when the server starts
var Subservice = map[string]subservice{
	"hus": {
		Domain: Domain{
			Name: "hus",
			URL:  "http://localhost:9090", // cloudhus.com
		},
		Subdomains: map[string]Domain{
			"auth": {
				Name: "auth",
				URL:  "http://localhost:9090", // auth.cloudhus.com
			},
			"api": {
				Name: "api",
				URL:  "http://localhost:9090", // api.cloudhus.com
			},
		},
	},
	"lifthus": {
		Domain: Domain{
			Name: "lifthus",
			URL:  "http://localhost:3000", // lifthus.com
		},
		Subdomains: map[string]Domain{
			"auth": {
				Name: "auth",
				URL:  "http://localhost:9091", // auth.lifthus.com
			},
			"api": {
				Name: "api",
				URL:  "http://localhost:9091", // api.lifthus.com
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
				URL:  "https://auth.surfhus.com", // auth.surfhus.com
			},
			"api": {
				Name: "api",
				URL:  "https://api.surfhus.com", // api.surfhus.com
			},
		},
	},
}
