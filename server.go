package main

import (
	"backend/database"      // Importing custom database package
	routes "backend/routes" // Importing custom routes package
	"context"               // Importing context package for handling contexts
	"os"                    // Importing os package for accessing environment variables

	"github.com/gin-contrib/cors"       // Importing CORS middleware package for Gin
	"github.com/gin-gonic/gin"          // Importing Gin web framework package
	"go.mongodb.org/mongo-driver/mongo" // Importing MongoDB driver for Go
)

var client *mongo.Client // Declaring a global variable to hold the MongoDB client instance

func mongoFun(c *gin.Context) {
	c.String(200, "hello") // Responding with "hello" and status code 200
}

func main() {
	router := gin.Default() // Creating a Gin router with default middleware

	client = database.DBinstance() // Initializing the MongoDB client

	defer client.Disconnect(context.Background()) // Deferring the disconnection of MongoDB client until the program exits
	// Defer ensures that the Disconnect method of the MongoDB client is called when the main function exits.
	// This helps to ensure proper cleanup and closing of resources when the program terminates.

	config := cors.DefaultConfig()                                  // Creating a default CORS configuration
	config.AllowOrigins = []string{"*"}                             // Allowing requests from any origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}  // Allowing specific HTTP methods
	config.AllowHeaders = []string{"Authorization", "Content-Type"} // Allowing specific headers
	router.Use(cors.New(config))                                    // Applying CORS middleware to the router
	// CORS middleware is used to enable Cross-Origin Resource Sharing, allowing resources on a web page to be requested from another domain.

	routes.AuthRoutes(router)  // Registering authentication routes
	routes.DiaryRoutes(router) // Registering diary routes

	router.Use(gin.Logger()) // Adding Gin's default logger middleware
	// Gin's logger middleware is used to log HTTP requests. It logs the method, path, and status code of each request.

	port := os.Getenv("PORT") // Getting the port number from the environment variable
	if port == "" {
		port = "8080" // Defaulting to port 8080 if not set
	}
	router.Run(":" + port) // Starting the HTTP server on the specified port
}
