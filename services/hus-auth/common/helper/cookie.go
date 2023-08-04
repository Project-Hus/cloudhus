package helper

import (
	"hus-auth/common/hus"
	"net/http"
)

// CookieMaker takes name and value and generates default lifthus auth cookie.
func AuthCookieMaker(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   hus.AuthCookieDomain,
		HttpOnly: true,
		Secure:   hus.CookieSecure,
		SameSite: hus.CookieSameSite,
	}
}

// HSTCookieMaker takes cookie's value and generate lifthus_st(which works like access token) cookie.
func HSTCookieMaker(value string) *http.Cookie {
	return AuthCookieMaker("hus_st", value)
}
