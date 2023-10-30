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

func TestTokenValid(t *testing.T) {
	// get token from login first
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2NjIzNzEsInVzZXJuYW1lIjoiYWRtaW4ifQ.owq1o80F7zIcA9f7WbCRA8TeqNbixwJTv94IUMgac40"
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	valid, _ := middlewares.AuthCheck(token, ctx)
	assert.Equal(t, true, valid)
}

func TestTokenInvalid(t *testing.T) {
	token := ""
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	valid, err := middlewares.AuthCheck(token, ctx)
	assert.Equal(t, false, valid)
	assert.Equal(t, err.Error(), "invalid token")
}

func TestTokenExpired(t *testing.T) {
	//set expired token in 15 second first(auth controller) and wait in 16 second
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2NDk5MDAsInVzZXJuYW1lIjoiYWRtaW4ifQ.osthIrLH-15cTgGi7roaPNRir0ad0hgek8h5VyqlfsI"
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	valid, err := middlewares.AuthCheck(token, ctx)

	assert.Equal(t, false, valid)
	assert.Equal(t, err.Error(), "expired token")
}

func TestShouldSuccessShowAllData(t *testing.T) {
	e := app.InitEcho()
	e.GET("/api/admin", ConnectController().FindAll)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))
	//get token from login first
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2ODg1NzUsInVzZXJuYW1lIjoiYWRtaW4ifQ.23hnrMTBrukzr8xqL-jz7ZiAHfQowgvui4i_vDYIiuY"

	req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
	req.Header.Set(echo.HeaderContentType, "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Unmarshal to map
	resBody := rec.Body.String()
	var resMap map[string]interface{}
	err := json.Unmarshal([]byte(resBody), &resMap)
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
	e := app.InitEcho()
	e.GET("/api/admin/:id", ConnectController().FindById)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))
	//get token from login first
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2ODg1NzUsInVzZXJuYW1lIjoiYWRtaW4ifQ.23hnrMTBrukzr8xqL-jz7ZiAHfQowgvui4i_vDYIiuY"

	id := "652d0ce5fc85ef0adf0ed692"
	req := httptest.NewRequest(http.MethodGet, "/api/admin/", nil)
	req.Header.Set(echo.HeaderContentType, "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

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
	e := app.InitEcho()
	e.POST("/api/admin", ConnectController().Create)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))
	requestBody := strings.NewReader(`{"username":"admin22","email":"admin22@gmail.com","password":"admin22"}`)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2ODg1NzUsInVzZXJuYW1lIjoiYWRtaW4ifQ.23hnrMTBrukzr8xqL-jz7ZiAHfQowgvui4i_vDYIiuY"
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/admin", requestBody)
	req.Header.Set(echo.HeaderContentType, "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
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
	e := app.InitEcho()
	e.PUT("/api/admin/:id", ConnectController().Update)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2ODg1NzUsInVzZXJuYW1lIjoiYWRtaW4ifQ.23hnrMTBrukzr8xqL-jz7ZiAHfQowgvui4i_vDYIiuY"

	requestBody := strings.NewReader(`{"username":"admin22","email":"admin22@gmail.com","password":"admin22"}`)
	reqCreate := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/admin", requestBody)
	reqCreate.Header.Set(echo.HeaderContentType, "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)

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
	reqUpdate.Header.Set("Authorization", "Bearer "+token)

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
	e := app.InitEcho()
	e.DELETE("/api/admin/:id", ConnectController().Delete)
	e.Use(middleware.KeyAuth(middlewares.AuthCheck))
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2ODg1NzUsInVzZXJuYW1lIjoiYWRtaW4ifQ.23hnrMTBrukzr8xqL-jz7ZiAHfQowgvui4i_vDYIiuY"

	requestBody := strings.NewReader(`{"username":"admindelete","email":"admindelete@gmail.com","password":"admin22"}`)
	reqCreate := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/admin", requestBody)
	reqCreate.Header.Set(echo.HeaderContentType, "application/json")
	reqCreate.Header.Set(echo.HeaderContentType, "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)

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
	req.Header.Set("Authorization", "Bearer "+token)
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
