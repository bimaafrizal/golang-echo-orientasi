package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthCheck(tokenString string, c echo.Context) (bool, error) {
	claims := TokenClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})

	if err != nil && parsed != nil {
		return false, errors.New("expired token")
	} else if err != nil && parsed == nil {
		return false, errors.New("invalid token")
	}

	c.Set("user_username", claims.Username)

	return true, nil
}
