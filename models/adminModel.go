package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID         primitive.ObjectID  `bson:"_id"`
	Name       *string              `json:"name" validate:"required" `
	Email	   *string				`json:"email" validate:"required"`
	Contact    *string				`json:"contact" validate:"required"`
	Password   *string				`json:"password" validate:"required"`
	Company		*string				`json:"company" `
	Token		*string				`json:"token" `
	Refresh_token *string 			`json:"refresh_token"`
	User_id		*string				`json:"user_id" `
	User_Type    *string				`json:"user_type"  validate:"required"`
}

type Client struct {
	ID           		primitive.ObjectID `bson:"_id"`
	Name         		*string            `json:"name" validate:"required"`
	Email        		*string            `json:"email" validate:"required"`
	Contact      		*string            `json:"contact" validate:"required"`
	Password     		*string            `json:"password" validate:"required"`
	Company      		*string            `json:"company" validate:"required"`
	Token        		*string            `json:"token"`
	Refresh_token 	*string            `json:"refresh_token"`
	User_id      		*string            `json:"user_id"`
}

type Request_to_admin struct {
	ID                  primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Title               *string `json:"name" validate:"required"`
	SendBy              *string `json:"send_by" validate:"required"`
	Discription         *string `json:"discription" validate:"required"`
	Short_discription   *string `json:"short_discription" validate:"required"`
	Sended_At           *string `json:"sended_at" validate:"required"`
	Status_review       *string `json:"status_review"`
	Status_reviewed_at  *string `json:"status_reviewed_at"`
}
