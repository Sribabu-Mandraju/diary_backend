package controllers

import (
	"backend/database"
	helper "backend/helpers"
	"backend/models"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

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
		defer cancel()
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

		if countByEmail > 0 {
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
		token := helper.GenerateAllTokens(
			*user.Email,
			*user.Name,
			*user.User_id,
		)
		user.Token = &token

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
		defer cancel()
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
				"error": "invalid credentials",
			})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "user not found",
			})
			return
		}

		token := helper.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.Name,
			*foundUser.User_id,
		)
		helper.UpdateAllTokens(token, *foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": *foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{
			"token":foundUser.Token,
		})
	}
}

func GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println("token string", tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		fmt.Println("token", token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		fmt.Println("claims", claims)

		userID, ok := claims["User_id"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}
		fmt.Println("id", userID)

		var user models.User
		ctx := context.TODO()
		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": user.User_id,
			"name":user.Name,
			"email":user.Email,
		})
	}
}
