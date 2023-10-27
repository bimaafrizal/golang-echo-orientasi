package test

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"orientasi-backend-golang/app"
	"orientasi-backend-golang/controller"
	"orientasi-backend-golang/middlewares"
	"orientasi-backend-golang/repository"
	"strings"
	"testing"
)

func ConnectController() controller.UserController {
	db := app.Init()
	validation := validator.New()
	userRepository := repository.NewUserRepository(db, validation)
	userController := controller.NewUserController(userRepository)

	return userController
}

func TestShouldSuccessShowAllData(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/admin", nil)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTgyOTkyNzEsInVzZXJuYW1lIjoiYWRtaW4ifQ.tfzKxc-ELgXmv85EW6l9JMnSDFSF-6ihHbMj_Nvl8II")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	err := ConnectController().FindAll(ctx)

	// Unmarshal to map
	resBody := rec.Body.String()
	var resMap map[string]interface{}
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		panic(err)
	}

	if assert.NoError(t, err) {
		message := resMap["message"].(string)
		assert.Equal(t, "Success Get All Data", message)
		assert.NotNil(t, resMap["data"])
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestShouldSuccessShowSpecificData(t *testing.T) {
	e := echo.New()
	id := "652d0ce5fc85ef0adf0ed692"

	req := httptest.NewRequest(http.MethodGet, "/api/admin/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/admin/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)

	err := ConnectController().FindById(ctx)

	// Unmarshal to map
	resBody := rec.Body.String()
	var resMap map[string]interface{}
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		panic(err)
	}

	if assert.NoError(t, err) {
		message := resMap["message"].(string)
		assert.Equal(t, "Success Get Specific Data", message)
		assert.NotNil(t, resMap["data"])
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestShouldSuccessCreateData(t *testing.T) {
	e := echo.New()
	requestBody := strings.NewReader(`{"username":"admin22","email":"admin22@gmail.com","password":"admin22"}`)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/admin", requestBody)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))

	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := ConnectController().Create(ctx)

	// Unmarshal to map
	resBody := rec.Body.String()
	var resMap map[string]interface{}
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		panic(err)
	}

	if assert.NoError(t, err) {
		message := resMap["message"].(string)
		assert.Equal(t, "Success Create Data", message)
		assert.NotNil(t, resMap["data"])
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestShouldSuccessUpdateData(t *testing.T) {
	e := echo.New()
	requestBody := strings.NewReader(`{"username":"admin22","email":"admin22@gmail.com","password":"admin22"}`)
	reqCreate := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/admin", requestBody)
	reqCreate.Header.Set(echo.HeaderContentType, "application/json")

	recCreate := httptest.NewRecorder()
	ctxCreate := e.NewContext(reqCreate, recCreate)
	err := ConnectController().Create(ctxCreate)

	resBody := recCreate.Body.String()

	// Unmarshal to map
	var resMap map[string]interface{}
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		panic(err)
	}

	// Get ID
	id := resMap["data"].(map[string]interface{})["id"].(string)

	requestUpdate := strings.NewReader(`{"username":"admin_update22","email":"admin22@gmail.com","password":"admin22"}`)
	reqUpdate := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/admin/652d0ce5fc85ef0adf0ed692", requestUpdate)
	reqUpdate.Header.Set(echo.HeaderContentType, "application/json")
	recUpdate := httptest.NewRecorder()
	ctxUpdate := e.NewContext(reqUpdate, recUpdate)

	ctxUpdate.SetPath("/api/admin/:id")
	ctxUpdate.SetParamNames("id")
	ctxUpdate.SetParamValues(id)

	err = ConnectController().Update(ctxUpdate)

	// Unmarshal to map
	resBodyAssert := recUpdate.Body.String()
	var resMapAssert map[string]interface{}
	err = json.Unmarshal([]byte(resBodyAssert), &resMapAssert)
	if err != nil {
		panic(err)
	}
	if assert.NoError(t, err) {
		message := resMapAssert["message"].(string)
		assert.Equal(t, "Success Update Data", message)
		assert.NotNil(t, resMapAssert["data"])
		assert.Equal(t, http.StatusOK, recUpdate.Code)
	}
}

func TestShouldSuccessDeleteData(t *testing.T) {
	e := echo.New()

	requestBody := strings.NewReader(`{"username":"admindelete","email":"admindelete@gmail.com","password":"admin22"}`)
	reqCreate := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/admin", requestBody)
	reqCreate.Header.Set(echo.HeaderContentType, "application/json")

	recCreate := httptest.NewRecorder()
	ctxCreate := e.NewContext(reqCreate, recCreate)
	err := ConnectController().Create(ctxCreate)

	resBody := recCreate.Body.String()

	// Unmarshal to map
	var resMap map[string]interface{}
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		panic(err)
	}

	// Get ID
	id := resMap["data"].(map[string]interface{})["id"].(string)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/admin", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/admin/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)

	err = ConnectController().Delete(ctx)
	// Unmarshal to map
	resBodyAssert := rec.Body.String()
	var resMapAssert map[string]interface{}
	err = json.Unmarshal([]byte(resBodyAssert), &resMapAssert)
	if err != nil {
		panic(err)
	}
	if assert.NoError(t, err) {
		message := resMapAssert["message"].(string)
		assert.Equal(t, "Success Delete Data", message)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
