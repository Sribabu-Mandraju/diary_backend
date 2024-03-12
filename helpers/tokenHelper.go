package helper

import (
	"context" // Importing context package for managing contexts
	"fmt"     // Importing fmt package for formatted I/O
	"log"     // Importing log package for logging
	"time"    // Importing time package for handling time-related operations

	"github.com/dgrijalva/jwt-go"               // Importing jwt-go package for JWT functionality
	"go.mongodb.org/mongo-driver/bson"          // Importing bson package for BSON document marshaling and unmarshaling
	"go.mongodb.org/mongo-driver/mongo"         // Importing mongo package for MongoDB operations
	"go.mongodb.org/mongo-driver/mongo/options" // Importing options package for MongoDB options

	"backend/database" // Importing database package for accessing MongoDB client and collections
)

// SignedDetails struct represents the details included in the JWT claims.
type SignedDetails struct {
	Email              string `json:"email"`    // Email field represents the email address of the user
	Name               string `json:"name"`     // Name field represents the name of the user
	Password           string `json:"password"` // Password field represents the password of the user
	User_id            string `json:"user_id"`  // User_id field represents the user ID
	jwt.StandardClaims        // Embedded field for standard JWT claims like expiration time
}

// userCollection represents the MongoDB collection for users.
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

// SECRET_KEY is the secret key used for JWT signing and validation.
var SECRET_KEY = "HELLO_WORLD"

// GenerateAllTokens function generates JWT token with specified claims.
func GenerateAllTokens(
	email string,
	name string,
	user_id string,
) (signedToken string) {
	// Creating JWT claims with specified details.
	claims := &SignedDetails{
		Email:   email,
		Name:    name,
		User_id: user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(), // Token expiration time set to 24 hours from current time
		},
	}

	// Signing the JWT token with the specified claims using HS256 algorithm and secret key.
	var err error
	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return signedToken
}

// UpdateAllTokens function updates the authentication token for a user in the database.
func UpdateAllTokens(signedToken string, userId string) {
	// Creating a context with a timeout for MongoDB operations.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Creating update operation document.
	updateObj := bson.D{
		{"$set", bson.D{
			{"token", signedToken},     // Setting the token field with the new token
			{"updated_at", time.Now()}, // Updating the updated_at field with current time
		}},
	}

	upsert := true                            // Indicates whether to perform an upsert operation
	filter := bson.M{"user_id": userId}       // Filter to identify the user document to update
	opt := options.Update().SetUpsert(upsert) // Setting options for update operation

	// Performing the update operation on the user collection.
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		updateObj,
		opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}
}

// ValidateToken function validates the JWT token and returns the claims if valid.
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	// Parsing and validating the JWT token with specified claims.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("HELLO_WORLD"), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token in invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
