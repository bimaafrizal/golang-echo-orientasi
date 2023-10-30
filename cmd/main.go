package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"orientasi-backend-golang/app"
	"orientasi-backend-golang/controller"
	"orientasi-backend-golang/middlewares"
	"orientasi-backend-golang/repository"
)

func main() {
	db := app.Init()
	validation := validator.New()
	userRepository := repository.NewUserRepository(db, validation)
	userController := controller.NewUserController(userRepository)
	authController := controller.NewAuthController(db, validation)

	E := app.InitEcho()

	E.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World Welcome To Orientasi Backend Golang")
	})

	// crud user
	crudGroup := E.Group("/api")
	crudGroup.Use(middleware.KeyAuth(middlewares.AuthCheck))
	crudGroup.GET("/admin", userController.FindAll)
	crudGroup.GET("/admin/:id", userController.FindById)
	crudGroup.DELETE("/admin/:id", userController.Delete)
	crudGroup.POST("/admin", userController.Create)
	crudGroup.PUT("/admin/:id", userController.Update)

	//auth
	E.POST("/api/login", authController.Login)

	err := E.Start(":8080")
	if err != nil {
		panic(err)
	}
}
