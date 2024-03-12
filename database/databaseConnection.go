package database

import (
	"context" // Importing context package for managing contexts
	"fmt"     // Importing fmt package for formatted I/O
	"log"     // Importing log package for logging
	"time"    // Importing time package for handling time-related operations

	"go.mongodb.org/mongo-driver/mongo"         // Importing mongo package for MongoDB operations
	"go.mongodb.org/mongo-driver/mongo/options" // Importing options package for MongoDB options
)

// DBinstance function creates and returns a MongoDB client instance.
func DBinstance() *mongo.Client {
	// MongoDB connection string.
	MongoDb := "mongodb+srv://sribabu:63037sribabu@atlascluster.k6u2oy9.mongodb.net/?retryWrites=true&w=majority"

	// Creating a new MongoDB client instance.
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err) // Exiting the program if client creation fails.
	}

	// Creating a context with a timeout for MongoDB connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Cancelling the context to release resources.

	// Establishing a connection to MongoDB.
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err) // Exiting the program if connection fails.
	}

	fmt.Println("Connected to MongoDB!") // Printing a success message upon successful connection.

	return client // Returning the MongoDB client instance.
}

// Client represents the MongoDB client instance.
var Client *mongo.Client = DBinstance()

// OpenCollection function establishes a connection with a collection in the database.
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Getting the specified collection from the database.
	var collection *mongo.Collection = client.Database("diary").Collection(collectionName)

	return collection // Returning the MongoDB collection instance.
}
