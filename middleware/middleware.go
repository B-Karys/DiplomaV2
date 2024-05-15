package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func LoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract token from header
		token := c.Request().Header.Get("Authorization")

		// Check if token is present
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing or invalid token"})
		}

		// (Optional) Validate token logic here

		// If token is valid, proceed with the next handler
		return next(c)
	}
}
