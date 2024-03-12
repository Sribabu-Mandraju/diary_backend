package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive" // Importing primitive package from MongoDB driver for BSON operations
)

// Diary struct represents the structure of a diary entry.
type Diary struct {
	ID          primitive.ObjectID `bson:"_id"`                       // ID field represents the unique identifier of the diary entry in BSON format
	Title       *string            `json:"title" validate:"required"` // Title field represents the title of the diary entry in JSON format
	Description *string            `json:"desc" validate:"required"`  // Description field represents the description of the diary entry in JSON format
	Date        *string            `json:"date" validate:"required"`  // Date field represents the date of the diary entry in JSON format
	User_id     *string            `json:"user_id" bson:"user_id"`    // User_id field represents the user ID associated with the diary entry in both JSON and BSON formats
}

// In Go, the `json` tag specifies the JSON key corresponding to the struct field when marshaling/unmarshaling JSON data.
// In this struct, `json:"title"` specifies that the Title field will be encoded/decoded as "title" in JSON.
// Similarly, `json:"desc"` specifies the JSON key for the Description field, `json:"date"` for the Date field, and `json:"user_id"` for the User_id field.

// In MongoDB, the `bson` tag specifies the BSON key corresponding to the struct field when working with MongoDB documents.
// In this struct, `bson:"_id"` specifies that the ID field will be represented as "_id" in BSON documents, and `bson:"user_id"` for the User_id field.

// The `validate` tag is not a standard Go tag but is often used with validation libraries like "github.com/go-playground/validator" to specify validation rules for struct fields.
// In this case, `validate:"required"` specifies that the Title, Description, and Date fields are required and must be provided when creating or updating a Diary object.
