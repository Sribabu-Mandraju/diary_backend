package database

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

//DBinstance func
func DBinstance() *mongo.Client {


    MongoDb := "mongodb+srv://sribabu:63037sribabu@atlascluster.k6u2oy9.mongodb.net/?retryWrites=true&w=majority"

    client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

    return client
}

//Client Database instance
var Client *mongo.Client = DBinstance()

//OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

    var collection *mongo.Collection = client.Database("u3technologies").Collection(collectionName)

    return collection
}