package auth

import (
	"fmt"
	"hus-auth/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// TokenEmbeddingHandler godoc
// @Router       /hus [post]
// @Summary      processes google auth and redirect with refresh token in url.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the refresh token will be expired in 7 days.
// @Tags         auth
// @Accept       json
// @Param        jwt body string true "Google ID token"
// @Success      301 "to /auth/{token_string}"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
func (ac authApiController) TokenEmbeddingHandler(c echo.Context) error {
	// get refresh token from header
	refreshToken := c.Request().Header.Get("Authorization")
	// validate refresh token
	_, err := service.ValidateRefreshToken(c.Request().Context(), ac.Client, refreshToken)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	// set cookie with refresh token
	cookie := &http.Cookie{
		Name:     "hus-refresh-token",
		Value:    refreshToken,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	fmt.Println("yo")
	// get whole set-cookie from echo context
	cookies := c.Cookies()

	for _, cookie := range cookies {
		fmt.Println(cookie.Name + cookie.Value)
	}

	return c.NoContent(http.StatusOK)
}

// RefreshTokenRequestHandler godoc
// @Router       /refresh [get]
// @Summary      gets refresh token in the header and returns access token in the cookie after validation.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the access token will be expired in 10 minutes.
// @Tags         auth
// @Param        jwt header string true "Refresh token"
// @Success      201 "Access token created"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
func (ac authApiController) RefreshTokenRequestHandler(c echo.Context) error {
	// get refresh token from cookie
	refreshToken, err := c.Cookie("hus-refresh-token")
	fmt.Println(refreshToken.Value)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return c.String(http.StatusOK, refreshToken.Value)
}

// AccessTokenRequestHandler godoc
// @Router       /access [get]
// @Summary      gets refresh token in the header and returns access token in the cookie after validation.
// @Description  validates the google ID token and redirects with hus refresh token to /auth/{token_string}.
// @Description the access token will be expired in 10 minutes.
// @Tags         auth
// @Param        jwt header string true "Refresh token"
// @Success      201 "Access token created"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal Server Error"
func (ac authApiController) AcessTokenRequestHandler(c echo.Context) error {
	// get refresh token from header
	refreshToken := c.Request().Header.Get("refresh_token")
	// validate refresh token
	refreshTokenValidated, err := service.ValidateRefreshToken(c.Request().Context(), ac.Client, refreshToken)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	// get user's uuid from refresh token
	uid := refreshTokenValidated["uid"].(string)

	// Create a new access token with 10 minutes expiration time.
	accessTokenSigned, err := service.GetNewAccessToken(c.Request().Context(), ac.Client, uid)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = accessTokenSigned
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.NoContent(http.StatusCreated)
}
