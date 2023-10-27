package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"os"
)

type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthCheck(tokenString string, c echo.Context) (bool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return false, err
	}
	envVariable := os.Getenv("JWT_SECRET_KEY")

	claims := TokenClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(envVariable), nil
	})
	if err != nil {
		return false, errors.New("invalid token")
	}
	if !parsed.Valid {
		return false, errors.New("expired token")
	}

	c.Set("user_username", claims.Username)

	return true, nil
}
