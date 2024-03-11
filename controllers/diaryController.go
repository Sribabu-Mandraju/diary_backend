package controllers

import (
	"backend/database"
	"backend/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var diaryCollection *mongo.Collection = database.OpenCollection(database.Client, "diary")

func GetAllDiaries() gin.HandlerFunc {
	return func(c *gin.Context) {
		var diaries []models.Diary

		// Extract user ID from request parameters
		userID := c.Param("userid")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Define filter to fetch diaries for the specified user ID
		filter := bson.M{"user_id": userID}

		cursor, err := diaryCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching diaries"})
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var diary models.Diary
			if err := cursor.Decode(&diary); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding diaries"})
				return
			}
			diaries = append(diaries, diary)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over diaries"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"diaries": diaries,
		})
	}
}

func GetDiaryById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var diary models.Diary

		id := c.Param("id")
		fmt.Println("id", id)
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = diaryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&diary)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
			return
		}

		c.JSON(http.StatusOK, diary)
	}
}

func CreateDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var diary models.Diary

		// fmt.Println("inside create dirary")

		if err := c.BindJSON(&diary); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationErr := validate.Struct(diary)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		diary.ID = primitive.NewObjectID()

		// fmt.Println("diary id", diary.ID)

		_, err := diaryCollection.InsertOne(ctx, diary)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create diary entry",
			})
			return
		}

		// fmt.Println("diary", diary)
		c.JSON(http.StatusCreated, gin.H{
			"message": "Diary entry created successfully",
			"diary":   diary,
		})

	}
}

func UpdateDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract diary ID from the request parameters
		id := c.Param("id") //dosucment id
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
			return
		}

		// Parse JSON data from the request body into an updated diary struct
		var updatedDiary models.Diary
		if err := c.BindJSON(&updatedDiary); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Set the ID field of the updated diary struct
		updatedDiary.ID = objID

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Update the diary entry in the database
		filter := bson.M{"_id": objID}
		update := bson.M{"$set": updatedDiary}
		_, err = diaryCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diary entry"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, gin.H{"message": "Diary updated successfully"})
	}
}

func DeleteDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract diary ID from the request parameters
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
			return
		}

		// Define a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Delete the diary entry from the database
		filter := bson.M{"_id": objID}
		_, err = diaryCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diary entry"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, gin.H{"message": "Diary deleted successfully"})
	}
}
