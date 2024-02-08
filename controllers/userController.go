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
	"github.com/dgrijalva/jwt-go"

)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var clientCollection2 *mongo.Collection = database.OpenCollection(database.Client, "clients")
var requestCollection *mongo.Collection = database.OpenCollection(database.Client, "requests")


var validate = validator.New()
var secretKey = []byte("HELLO_WORLD")






func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}


func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := err == nil
	msg := ""
	if !check {
		msg = fmt.Sprintf("email or password not matched")
	}
	return check, msg
}


func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		countByEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred",
			})
			return
		}

		countByContact, err := userCollection.CountDocuments(ctx, bson.M{"contact": user.Contact})
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

		password := HashPassword(*user.Password)
		user.Password = &password

		user.ID = primitive.NewObjectID()
		userID := user.ID.Hex()
		user.User_id = &userID
		token, refreshToken := helper.GenerateAllTokens(
			*user.Email,
			*user.Name,
			*user.Company,
			*user.User_id,
			*user.Contact,
		)
		user.Token = &token
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user":    user,
		})

		return
	}
}


func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "email or password not matched",
			})
			return
		}

		passwordIsValid, _ := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "hello world",
			})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "HELLO",
			})
			return
		}

		token, refreshToken := helper.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.Name,
			*foundUser.Company,
			*foundUser.User_id,
			*foundUser.Contact,
		)
		helper.UpdateAllTokens(token, refreshToken, *foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": *foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetAllAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		cursor, err := userCollection.Find(ctx, bson.M{"user_type":"ADMIN"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
	
		var users []bson.M
		err = cursor.All(ctx, &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, users)
	}
}

func GetAdminByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		userIDParam := c.Param("id")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	
		var user bson.M
		err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	
		c.JSON(http.StatusOK, user)
	}
}




func GetUserInfo() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println("token string",tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			return
		}
	
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		fmt.Println("token",token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
	
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		fmt.Println("claims",claims)
	
		userID, ok := claims["User_id"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}
		fmt.Println("id",userID)
	
		var user models.User
		ctx := context.TODO()
		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}
	
		c.JSON(http.StatusOK, user)
	}
	
}

func GetAllRequests() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Fetch all request elements from the database
        var requests []models.Request_to_admin
        cursor, err := requestCollection.Find(context.Background(), bson.M{})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to fetch requests",
            })
            return
        }
        defer cursor.Close(context.Background())

        if err := cursor.All(context.Background(), &requests); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to decode requests",
            })
            return
        }

        // Return the list of request elements to the client
        c.JSON(http.StatusOK, requests)
    }
}



func ApproveOrRejectRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the request ID from the URL path parameters
        requestID := c.Param("id")

        // Parse requestID to ObjectId (assuming MongoDB ObjectId)
        objectID, err := primitive.ObjectIDFromHex(requestID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "invalid request ID",
            })
            return
        }

        // Get the action (approve or reject) from the request body
        var action struct {
            Action string `json:"action" binding:"required"`
        }
        if err := c.BindJSON(&action); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "invalid request body",
            })
            return
        }

        // Update the status of the request in the database
        update := bson.M{"status_review": action.Action}
        _, err = requestCollection.UpdateOne(
            context.Background(),
            bson.M{"_id": objectID},
            bson.M{"$set": update},
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to update request",
            })
            return
        }

        // Respond with a success message
        c.JSON(http.StatusOK, gin.H{
            "msg": fmt.Sprintf("request %s successfully", action.Action),
        })
    }
}


func GetAllClients() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		cursor, err := clientCollection2.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
	
		var users []bson.M
		err = cursor.All(ctx, &users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, users)
	}
}

func GetClientByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
	
		userIDParam := c.Param("id")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	
		var user bson.M
		err = clientCollection2.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	
		c.JSON(http.StatusOK, user)
	}
}


