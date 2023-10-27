package controller

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"orientasi-backend-golang/helper"
	"orientasi-backend-golang/repository"
)

type UserControllerImpl struct {
	UserRepository repository.UserRepository
}

func NewUserController(userRepository repository.UserRepository) UserController {
	return &UserControllerImpl{
		UserRepository: userRepository,
	}
}

func (controller *UserControllerImpl) FindAll(c echo.Context) error {
	user, err := controller.UserRepository.FindAll()
	if err != nil {
		return helper.InternalServerErrorResponse()
	}

	response := helper.Response{Code: 200, Message: "Success Get All Data", Data: user}
	return c.JSON(200, response)
}

func (controller *UserControllerImpl) FindById(c echo.Context) error {
	userId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return helper.BadRequestResponse(err)
	}

	user, err := controller.UserRepository.FindById(objectId)
	if err != nil {
		return helper.InternalServerErrorResponse()
	}

	response := helper.Response{Code: 200, Message: "Success Get Specific Data", Data: user}
	return c.JSON(200, response)
}

func (controller *UserControllerImpl) Create(c echo.Context) error {
	user := new(helper.CreateUserRequest)
	err := c.Bind(user)
	if err != nil {
		return helper.BadRequestResponse(err)
	}
	uniqueEmail := controller.UserRepository.CheckUniqueEmail(user.Email, primitive.NilObjectID)
	uniqueUsername := controller.UserRepository.CheckUniqueUsername(user.Username, primitive.NilObjectID)

	if uniqueUsername == false {
		err := "Username"
		return helper.UniqueRequestResponse(err)
	}
	if uniqueEmail == false {
		err := "Email"
		return helper.UniqueRequestResponse(err)
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return helper.InternalServerErrorResponse()
	}

	userUpdateHash := helper.CreateUserRequest{
		Username: user.Username,
		Password: hashPassword,
		Email:    user.Email,
	}

	err, id := controller.UserRepository.Create(&userUpdateHash)
	dataResponse := helper.ResponseWithId{Id: id.Hex(), Username: user.Username, Email: user.Email, Password: userUpdateHash.Password}

	if err != nil {
		return helper.BadRequestResponse(err)
	}

	response := helper.Response{Code: http.StatusCreated, Message: "Success Create Data", Data: &dataResponse}
	return c.JSON(http.StatusCreated, response)
}

func (controller *UserControllerImpl) Update(c echo.Context) error {
	userId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return helper.BadRequestResponse(err)
	}

	_, err = controller.UserRepository.FindById(objectId)
	if err != nil {
		return helper.BadRequestResponse(err)
	}

	user := new(helper.UpdateUserRequest)
	err = c.Bind(user)
	if err != nil {
		return helper.BadRequestResponse(err)
	}
	uniqueEmail := controller.UserRepository.CheckUniqueEmail(user.Email, objectId)
	uniqueUsername := controller.UserRepository.CheckUniqueUsername(user.Username, objectId)

	if uniqueUsername == false {
		err := "Username"
		return helper.UniqueRequestResponse(err)
	}

	if uniqueEmail == false {
		err := "Email"
		return helper.UniqueRequestResponse(err)
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return helper.InternalServerErrorResponse()
	}
	userUpdateHash := helper.UpdateUserRequest{
		Username: user.Username,
		Password: hashPassword,
		Email:    user.Email,
	}
	update, err := controller.UserRepository.Update(&userUpdateHash, objectId)
	if err != nil {
		return err
	}

	response := helper.Response{Code: http.StatusOK, Message: "Success Update Data", Data: update}
	return c.JSON(http.StatusOK, response)
}

func (controller *UserControllerImpl) Delete(c echo.Context) error {
	userId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return helper.BadRequestResponse(err)
	}
	err = controller.UserRepository.Delete(objectId)
	if err != nil {
		return helper.BadRequestResponse(err)
	}

	response := helper.Response{Code: http.StatusOK, Message: "Success Delete Data"}
	return c.JSON(http.StatusOK, response)
}
