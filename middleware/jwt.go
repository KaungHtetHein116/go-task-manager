package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "Unthorized",
				})
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "Invalid token",
				})
			}

			// Add user id to context
			claims := token.Claims.(jwt.MapClaims)
			c.Set("user_id", uint(claims["user_id"].(float64)))

			return next(c)
		}
	}
}
