package middleware

import (
	"github.com/labstack/echo/v4"
)

var CorsOrigins = map[string]string{}

// SetHusCorsHeaders sets headers for CORS.
func SetHusCorsHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// If your Backend is deployed in AWS and using API Gateway to call through,
		//then all these headers need to be applied in API Gateway level also.
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		// Access-Control-Allow-Methodsand Access-Control-Allow-Headersshould contain the same value
		//as requested in Access-Control-request-Methodsand Access-Control-request-Headersrespectively.
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, *")
		return next(c)
	}
}
