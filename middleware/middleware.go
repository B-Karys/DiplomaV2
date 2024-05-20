package middleware2

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// JWTSecretKey is the key used to sign the JWT token.
// Replace this with your actual secret key.
var JWTSecretKey = []byte(os.Getenv("JWT_SECRET"))

func LoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if user is already in context
		if user := c.Get("user"); user != nil {
			return next(c)
		}

		// Retrieve token from cookie
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing or invalid cookie"})
		}
		tokenString := cookie.Value

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}
			return JWTSecretKey, nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		// Check if token is valid and get the claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check token expiration
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token has expired"})
			}

			// Extract subject (user ID)
			var subject string
			switch sub := claims["sub"].(type) {
			case string:
				subject = sub
			case float64:
				subject = fmt.Sprintf("%.0f", sub)
			default:
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token subject"})
			}

			// If token is valid, proceed with the next handler
			c.Set("userID", subject)
			return next(c)
		} else {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}
	}
}
