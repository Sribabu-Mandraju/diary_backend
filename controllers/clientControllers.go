package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
	"backend/models"
	"backend/database"
	helper "backend/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var clientCollection *mongo.Collection = database.OpenCollection(database.Client, "clients")
var validateClient = validator.New()

func HashPasswordClient(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPasswordClient(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := err == nil
	msg := ""
	if !check {
		msg = fmt.Sprintf("email or password not matched")
	}
	return check, msg
}


func ClientRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var client models.Client

		if err := c.BindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationErr := validateClient.Struct(client)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		countByEmail, err := clientCollection.CountDocuments(ctx, bson.M{"email": client.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred",
			})
			return
		}

		countByContact, err := clientCollection.CountDocuments(ctx, bson.M{"contact": client.Contact})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred",
			})
			return
		}

		// Check if either email or contact already exists
		if countByEmail > 0 && countByContact > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "this email or contact already exists",
			})
			return
		}

		password := HashPasswordClient(*client.Password)
		client.Password = &password

		client.ID = primitive.NewObjectID()
		clientID := client.ID.Hex()
		client.User_id = &clientID

		if client.Email == nil || client.Name == nil || client.Company == nil || client.Password == nil || client.Contact == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing required fields",
			})
			return
		}

		token, refreshToken := helper.GenerateAllTokens(
			*client.Email,
			*client.Name,
			*client.Company,
			*client.User_id,
			*client.Contact,
		)
		client.Token = &token
		client.Refresh_token = &refreshToken

		_, insertErr := clientCollection.InsertOne(ctx, client)
		if insertErr != nil {
			msg := fmt.Sprintf("client item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "client registered successfully",
			"client":  client,
		})

		return
	}
}



func ClientLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var client models.Client
		var foundClient models.Client

		if err := c.BindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := clientCollection.FindOne(ctx, bson.M{"email": client.Email}).Decode(&foundClient)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "email or password not matched",
			})
			return
		}

		passwordIsValid, msg := VerifyPasswordClient(*client.Password, *foundClient.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": msg,
			})
			return
		}

		if foundClient.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "HELLO",
			})
			return
		}

		token, refreshToken := helper.GenerateAllTokens(
			*foundClient.Email,
			*foundClient.Name,
			*foundClient.Company,
			*foundClient.User_id,
			*foundClient.Contact,
		)
		helper.UpdateAllTokens(token, refreshToken, *foundClient.User_id)
		err = clientCollection.FindOne(ctx, bson.M{"user_id": *foundClient.User_id}).Decode(&foundClient)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"client":  foundClient,
		})
	}
}

func SendRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.Request_to_admin

		// Parse the incoming JSON request body into the Requests struct
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get the current time and format it as desired
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		request.Sended_At = &currentTime

		// Insert the request data into the MongoDB collection
		result, err := requestCollection.InsertOne(context.Background(), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Respond with success message and inserted data
		c.JSON(http.StatusOK, gin.H{
			"msg":  "successfully request sent",
			"data": request,
			"status":result,
		})
	}
}



