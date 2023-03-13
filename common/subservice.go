package common

type subservice struct {
	Name       string          `json:"name"`
	Prtc       string          `json:"prtc"`
	Domain     string          `json:"domain"`
	Subdomains map[string]bool `json:"subdomains"`
}

// later it would be good to have this info in a database and to initialize this map when the server starts
const httpPrtc = "http"

var Subservice = map[string]subservice{
	"hus": {
		Name:   "cloudhus",
		Prtc:   httpPrtc,
		Domain: "cloudhus.com",
		Subdomains: map[string]bool{
			"auth": true, // auth.cloudhus.com
			"api":  true, // api.cloudhus.com
		},
	},
	"lifthus": {
		Name:   "lifthus",
		Prtc:   httpPrtc,
		Domain: "lifthus.com",
		Subdomains: map[string]bool{
			"auth": true, // auth.lifthus.com
			"api":  true, // api.lifthus.com
		},
	},
	"surfhus": {
		Name:   "surfhus",
		Prtc:   httpPrtc,
		Domain: "surfhus.com",
		Subdomains: map[string]bool{
			"auth": true, // auth.surfhus.com
			"api":  true,
		},
	},
}
