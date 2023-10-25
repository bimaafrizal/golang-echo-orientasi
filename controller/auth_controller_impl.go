package controller

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"orientasi-backend-golang/helper"
	"orientasi-backend-golang/model"
	"os"
	"time"
)

type AuthControllerImpl struct {
	db        *mongo.Database
	validator *validator.Validate
}

func NewAuthController(db *mongo.Database, v *validator.Validate) AuthController {
	return &AuthControllerImpl{db: db, validator: v}
}

func (a *AuthControllerImpl) Login(c echo.Context) error {
	userRequest := new(helper.LoginRequest)
	err := c.Bind(userRequest)
	err = a.validator.Struct(userRequest)
	if err != nil {
		return helper.BadRequestResponse(err)
	}

	filter := bson.M{"username": userRequest.Username}
	var user model.User
	err = a.db.Collection(model.TableUsers).FindOne(context.TODO(), filter).Decode(&user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}
	token, err := GenerateToken(userRequest.Username)
	if err != nil {
		return helper.InternalServerErrorResponse()
	}

	return c.JSON(http.StatusOK, map[string]string{
		"code":    "200",
		"message": "success login",
		"token":   token,
	})
}

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	envVariable := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(envVariable))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
