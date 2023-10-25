package helper

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func InternalServerErrorResponse() error {
	errResponse := ResponseError{Code: 500, Message: "Internal Server Error"}
	return echo.NewHTTPError(http.StatusInternalServerError, errResponse)
}

func BadRequestResponse(err error) error {
	errResponse := ResponseError{Code: http.StatusBadRequest, Message: err.Error()}
	return echo.NewHTTPError(http.StatusBadRequest, errResponse)
}

func UniqueRequestResponse(err string) error {
	errResponse := ResponseError{Code: http.StatusConflict, Message: err + " Already Exist"}
	return echo.NewHTTPError(http.StatusConflict, errResponse)
}
