package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID         primitive.ObjectID  `bson:"_id"`
	Name       *string              `json:"name" validate:"required,min=5" `
	Email	   *string				`json:"email" validate:"required"`
	Password   *string				`json:"password" validate:"required"`
	Token		*string				`json:"token" `
	User_id		*string				`json:"user_id" `
}

