package repository

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"orientasi-backend-golang/helper"
	"orientasi-backend-golang/model"
)

type UserRepositoryImpl struct {
	db        *mongo.Database
	validator *validator.Validate
}

func NewUserRepository(db *mongo.Database, v *validator.Validate) UserRepository {
	return &UserRepositoryImpl{db: db, validator: v}
}

func (repository *UserRepositoryImpl) FindAll() ([]model.User, error) {
	result := []model.User{}
	ctx := context.TODO()

	query, err := repository.db.Collection(model.TableUsers).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func(query *mongo.Cursor, ctx context.Context) {
		err := query.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(query, ctx)

	for query.Next(ctx) {
		var row model.User
		err := query.Decode(&row)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (repository *UserRepositoryImpl) FindById(id primitive.ObjectID) (*model.User, error) {
	var users model.User
	filter := bson.D{{"_id", id}}
	err := repository.db.Collection(model.TableUsers).FindOne(context.Background(), filter).Decode(&users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (repository *UserRepositoryImpl) Create(user *helper.CreateUserRequest) (error, primitive.ObjectID) {
	err := repository.validator.Struct(user)
	if err != nil {
		return err, primitive.NilObjectID
	}
	data, err := repository.db.Collection(model.TableUsers).InsertOne(context.Background(), user)
	if err != nil {
		return err, primitive.NilObjectID
	}

	id := data.InsertedID.(primitive.ObjectID)
	return nil, id
}

func (repository *UserRepositoryImpl) Update(user *helper.UpdateUserRequest, id primitive.ObjectID) (*helper.UpdateUserRequest, error) {
	filter := bson.D{{"_id", id}}
	_, err := repository.db.Collection(model.TableUsers).ReplaceOne(context.Background(), filter, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repository *UserRepositoryImpl) Delete(id primitive.ObjectID) error {
	_, err := repository.db.Collection(model.TableUsers).DeleteOne(context.Background(), bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}
	return nil
}

// CheckUniqueEmail for validation unique email
func (repository *UserRepositoryImpl) CheckUniqueEmail(email string, idExcept primitive.ObjectID) bool {
	filter := bson.M{"email": email}

	if idExcept == primitive.NilObjectID {
		cur, err := repository.db.Collection(model.TableUsers).Find(context.Background(), filter)
		if err != nil {
			return false
		}
		return repository.MoreThanOne(cur)
	}

	filter = bson.M{"email": email, "_id": bson.M{"$ne": idExcept}}
	cur, err := repository.db.Collection(model.TableUsers).Find(context.Background(), filter)
	if err != nil {
		return false
	}
	return repository.MoreThanOne(cur)
}

// for validation unique username
func (repository *UserRepositoryImpl) CheckUniqueUsername(username string, idExcept primitive.ObjectID) bool {
	filter := bson.M{"username": username}

	if idExcept == primitive.NilObjectID {
		cur, err := repository.db.Collection(model.TableUsers).Find(context.Background(), filter)
		if err != nil {
			return false
		}
		return repository.MoreThanOne(cur)
	}

	filter = bson.M{"username": username, "_id": bson.M{"$ne": idExcept}}
	cur, err := repository.db.Collection(model.TableUsers).Find(context.Background(), filter)
	if err != nil {
		return false
	}
	return repository.MoreThanOne(cur)
}

func (repository *UserRepositoryImpl) MoreThanOne(cur *mongo.Cursor) bool {
	var items []model.User
	for cur.Next(context.TODO()) {
		var item model.User
		if err := cur.Decode(&item); err != nil {
			return false
		}
		items = append(items, item)
	}
	fmt.Println(len(items))
	if len(items) > 0 {
		return false
	}
	return true
}
