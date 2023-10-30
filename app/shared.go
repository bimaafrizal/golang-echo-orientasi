package app

import "github.com/labstack/echo/v4"

var E *echo.Echo

func InitEcho() *echo.Echo {
	E = echo.New()
	return E
}
