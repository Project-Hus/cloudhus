package auth

import (
	"fmt"
	"hus-auth/helper"
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
	_, err := helper.ValidateRefreshToken(c.Request().Context(), ac.Client, refreshToken)
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
