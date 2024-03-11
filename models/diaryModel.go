package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diary struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       *string            `json:"title" validate:"required"`
	Description *string            `json:"desc" validate:"required"`
	Date        *string            `json:"date" validate:"required"`
	User_id     *string            `json:"user_id" bson:"user_id"`
}
