package common

type subservice struct {
	Name       string            `json:"name"`
	Domain     string            `json:"domain"`
	Subdomains map[string]string `json:"subdomains"`
	URL        string            `json:"url"`
}

// later it would be good to have this info in a database and to initialize this map when the server starts
var Subservice = map[string]subservice{
	"hus": {
		Name:   "cloudhus",
		Domain: "cloudhus.com",
		Subdomains: map[string]string{
			"auth": "localhost:9090", // auth.cloudhus.com
			"api":  "localhost:9090", // api.cloudhus.com
		},
		URL: "http://localhost:9090", // https://cloudhus.com
	},
	"lifthus": {
		Name:   "lifthus",
		Domain: "lifthus.com",
		Subdomains: map[string]string{
			"auth": "auth.localhost:9091", // auth.lifthus.com
			"api":  "localhost:9091",      // api.lifthus.com
		},
		URL: "http://localhost:9091", // https://lifthus.com
	},
	"surfhus": {
		Name:   "surfhus",
		Domain: "surfhus.com",
		Subdomains: map[string]string{
			"auth": "auth.surfhus.com", // auth.surfhus.com
			"api":  "https://api.surfhus.com",
		},
		URL: "https://surfhus.com",
	},
}
