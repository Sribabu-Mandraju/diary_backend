package controllers

import (
	"backend/database" // Importing database package for accessing MongoDB collections
	"backend/models"   // Importing models package for defining data models
	"context"          // Importing context package for managing contexts
	// Importing fmt package for formatted I/O
	"net/http" // Importing net/http package for HTTP status codes
	"time"     // Importing time package for handling time-related operations

	"github.com/gin-gonic/gin"                   // Importing Gin web framework package
	"go.mongodb.org/mongo-driver/bson"           // Importing bson package for BSON document marshaling and unmarshaling
	"go.mongodb.org/mongo-driver/bson/primitive" // Importing primitive package for MongoDB object IDs
	"go.mongodb.org/mongo-driver/mongo"          // Importing mongo package for MongoDB operations
)

// diaryCollection represents the MongoDB collection for diaries.
var diaryCollection *mongo.Collection = database.OpenCollection(database.Client, "diary")

// GetAllDiaries function retrieves all diaries associated with a specific user.
func GetAllDiaries() gin.HandlerFunc {
	// Returning a Gin handler function to get all diaries.
	return func(c *gin.Context) {
		var diaries []models.Diary // Slice to hold retrieved diaries

		// Extracting user ID from request parameters
		userID := c.Param("userid")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Defining filter to fetch diaries for the specified user ID
		filter := bson.M{"user_id": userID}

		// Fetching diaries from the database based on the filter
		cursor, err := diaryCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching diaries"})
			return
		}
		defer cursor.Close(ctx)

		// Iterating over the retrieved diaries and appending them to the slice
		for cursor.Next(ctx) {
			var diary models.Diary
			if err := cursor.Decode(&diary); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding diaries"})
				return
			}
			diaries = append(diaries, diary)
		}

		// Checking for errors during cursor iteration
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over diaries"})
			return
		}

		// Responding with the retrieved diaries
		c.JSON(http.StatusOK, gin.H{
			"diaries": diaries,
		})
	}
}

// GetDiaryById function retrieves a diary entry by its ID.
func GetDiaryById() gin.HandlerFunc {
	// Returning a Gin handler function to get a diary by ID.
	return func(c *gin.Context) {
		var diary models.Diary // Variable to hold retrieved diary

		// Extracting diary ID from the request parameters
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Fetching diary from the database based on the ID
		err = diaryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&diary)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
			return
		}

		// Responding with the retrieved diary
		c.JSON(http.StatusOK, diary)
	}
}

// CreateDiary function creates a new diary entry.
func CreateDiary() gin.HandlerFunc {
	// Returning a Gin handler function to create a diary entry.
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var diary models.Diary

		// Parsing request body into diary struct
		if err := c.BindJSON(&diary); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validating diary struct fields
		validationErr := validate.Struct(diary)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		// Generating a new ObjectID for the diary
		diary.ID = primitive.NewObjectID()

		// Inserting the diary entry into the database
		_, err := diaryCollection.InsertOne(ctx, diary)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create diary entry",
			})
			return
		}

		// Responding with success message and created diary entry
		c.JSON(http.StatusCreated, gin.H{
			"message": "Diary entry created successfully",
			"diary":   diary,
		})

	}
}

// UpdateDiary function updates an existing diary entry.
func UpdateDiary() gin.HandlerFunc {
	// Returning a Gin handler function to update a diary entry.
	return func(c *gin.Context) {
		// Extracting diary ID from the request parameters
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
			return
		}

		// Parsing JSON data from the request body into an updated diary struct
		var updatedDiary models.Diary
		if err := c.BindJSON(&updatedDiary); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Setting the ID field of the updated diary struct
		updatedDiary.ID = objID

		// Defining a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Updating the diary entry in the database
		filter := bson.M{"_id": objID}
		update := bson.M{"$set": updatedDiary}
		_, err = diaryCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diary entry"})
			return
		}

		// Responding with success message
		c.JSON(http.StatusOK, gin.H{"message": "Diary updated successfully"})
	}
}

// DeleteDiary function deletes a diary entry.
func DeleteDiary() gin.HandlerFunc {
	// Returning a Gin handler function to delete a diary entry.
	return func(c *gin.Context) {
		// Extracting diary ID from the request parameters
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
			return
		}

		// Defining a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(),

			5*time.Second)
		defer cancel()

		// Deleting the diary entry from the database
		filter := bson.M{"_id": objID}
		_, err = diaryCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diary entry"})
			return
		}

		// Responding with success message
		c.JSON(http.StatusOK, gin.H{"message": "Diary deleted successfully"})
	}
}
