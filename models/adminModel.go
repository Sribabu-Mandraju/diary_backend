package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive" // Importing primitive package from MongoDB driver for BSON operations
)

// User struct represents the structure of a user.
type User struct {
	ID       primitive.ObjectID `bson:"_id"`                            // ID field represents the unique identifier of the user in BSON format
	Name     *string            `json:"name" validate:"required,min=5"` // Name field represents the name of the user in JSON format and must be at least 5 characters long
	Email    *string            `json:"email" validate:"required"`      // Email field represents the email address of the user in JSON format and is required
	Password *string            `json:"password" validate:"required"`   // Password field represents the password of the user in JSON format and is required
	Token    *string            `json:"token"`                          // Token field represents the authentication token associated with the user in JSON format
	User_id  *string            `json:"user_id"`                        // User_id field represents the user ID in JSON format
}

// In Go, the `json` tag specifies the JSON key corresponding to the struct field when marshaling/unmarshaling JSON data.
// In this struct, `json:"name"` specifies that the Name field will be encoded/decoded as "name" in JSON.
// Similarly, `json:"email"` specifies the JSON key for the Email field, `json:"password"` for the Password field, `json:"token"` for the Token field, and `json:"user_id"` for the User_id field.

// In MongoDB, the `bson` tag specifies the BSON key corresponding to the struct field when working with MongoDB documents.
// In this struct, `bson:"_id"` specifies that the ID field will be represented as "_id" in BSON documents.

// The `validate` tag is not a standard Go tag but is often used with validation libraries like "github.com/go-playground/validator" to specify validation rules for struct fields.
// In this case, `validate:"required,min=5"` specifies that the Name field is required and must be at least 5 characters long, `validate:"required"` specifies that the Email and Password fields are required.
