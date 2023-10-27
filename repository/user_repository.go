package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"orientasi-backend-golang/helper"
	"orientasi-backend-golang/model"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	FindById(id primitive.ObjectID) (*model.User, error)
	Create(user *helper.CreateUserRequest) (error, primitive.ObjectID)
	Update(user *helper.UpdateUserRequest, id primitive.ObjectID) (*helper.UpdateUserRequest, error)
	Delete(id primitive.ObjectID) error
	CheckUniqueEmail(email string, except primitive.ObjectID) bool
	CheckUniqueUsername(email string, except primitive.ObjectID) bool
}
