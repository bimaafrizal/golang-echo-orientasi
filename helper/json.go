package helper

type Response struct {
	Code    int         `bson:"code" json:"code"`
	Message string      `bson:"message" json:"message"`
	Data    interface{} `bson:"data" json:"data"`
}

type ResponseError struct {
	Code    int    `bson:"code" json:"code"`
	Message string `bson:"message" json:"message"`
}
