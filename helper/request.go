package helper

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=255" bson:"username"`
	Email    string `json:"email" validate:"required,min=5,email,max=255" bson:"email"`
	Password string `json:"password" validate:"required,min=5,max=255" bson:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=255" bson:"username"`
	Email    string `json:"email" validate:"required,min=5,email,max=255" bson:"email"`
	Password string `json:"password" validate:"required,min=5,max=255" bson:"password"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=255" bson:"username"`
	Password string `json:"password" validate:"required,min=5,max=255" bson:"password"`
}
