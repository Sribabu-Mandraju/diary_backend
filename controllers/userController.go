package controllers

import (
	"context"  // Importing context package for managing contexts
	"fmt"      // Importing fmt package for formatted I/O
	"log"      // Importing log package for logging
	"net/http" // Importing net/http package for HTTP status codes
	"time"     // Importing time package for handling time-related operations

	"backend/database"       // Importing database package for accessing MongoDB collections
	helper "backend/helpers" // Importing helper package for token generation and validation
	"backend/models"         // Importing models package for defining data models

	"github.com/gin-gonic/gin"                   // Importing Gin web framework package
	"github.com/go-playground/validator/v10"     // Importing validator package for struct field validation
	"go.mongodb.org/mongo-driver/bson"           // Importing bson package for BSON document marshaling and unmarshaling
	"go.mongodb.org/mongo-driver/bson/primitive" // Importing primitive package for MongoDB object IDs
	"go.mongodb.org/mongo-driver/mongo"          // Importing mongo package for MongoDB operations
	"golang.org/x/crypto/bcrypt"                 // Importing bcrypt package for password hashing and comparison
)

// userCollection represents the MongoDB collection for users.
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

// validate is a validator instance for struct field validation.
var validate = validator.New()

// HashPassword function hashes the provided password using bcrypt.
func HashPassword(password string) string {
	// Hashing the password using bcrypt with a cost factor of 14.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err) // Logging fatal error if password hashing fails.
	}
	return string(bytes) // Returning the hashed password.
}

// VerifyPassword function verifies the provided password against the user's stored hashed password.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// Comparing the provided password with the stored hashed password using bcrypt.
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := err == nil // Checking if password comparison was successful.
	msg := ""
	if !check {
		msg = fmt.Sprintf("email or password not matched")
	}
	return check, msg // Returning the result of password comparison and an error message if any.
}

// Register function handles user registration requests.
func Register() gin.HandlerFunc {
	// Returning a Gin handler function for user registration.
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		// Parsing request body into user struct.
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validating user struct fields.
		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		// Checking if user with the same email already exists.
		countByEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred",
			})
			return
		}

		if countByEmail > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "this email or contact already exists",
			})
			return
		}

		// Hashing the password before storing it.
		password := HashPassword(*user.Password)
		user.Password = &password

		// Generating unique user ID and assigning it to the user.
		user.ID = primitive.NewObjectID()
		userID := user.ID.Hex()
		user.User_id = &userID

		// Generating authentication token for the user.
		token := helper.GenerateAllTokens(
			*user.Email,
			*user.Name,
			*user.User_id,
		)
		user.Token = &token

		// Inserting the user data into the database.
		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		defer cancel()
		// Responding with success message and user data.
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user":    user,
		})

		return
	}
}

// Login function handles user login requests.
func Login() gin.HandlerFunc {
	// Returning a Gin handler function for user login.
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		// Parsing request body into user struct.
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Finding user by email in the database.
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		// Handling case where user is not found.
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "email or password not matched",
			})
			return
		}

		// Verifying the provided password with the stored hashed password.
		passwordIsValid, _ := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		// Handling case where password verification fails.
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "hello world",
			})
			return
		}

		// Generating authentication token for the user.
		token := helper.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.Name,
			*foundUser.User_id,
		)
		// Updating user's authentication token in the database.
		helper.UpdateAllTokens(token, *foundUser.User_id)
		// Retrieving updated user data from the database.
		err = userCollection.FindOne(ctx, bson.M{"user_id": *foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Responding with user data.
		c.JSON(http.StatusOK, foundUser)
	}
}
