package main

import (
	"context"
	"os"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"backend/database"
	"github.com/gin-contrib/cors"
	routes "backend/routes"
)

var client *mongo.Client


func mongoFun(c *gin.Context){
	c.String(200,"hello")
}



func main() {
	router := gin.Default()
	client = database.DBinstance()
	defer client.Disconnect(context.Background())
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Add your React app's origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

	routes.AuthRoutes(router)
	routes.ClientAuthRoutes(router)
	routes.UserRoutes(router)
	router.Use(gin.Logger())


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// package main

// import (
// 	"github.com/gin-gonic/gin"
//     "go.mongodb.org/mongo-driver/mongo"
// 	"modfile/handlers"
// 	"modfile/db"
// 	"github.com/gin-contrib/cors"
// )	
// var client *mongo.Client


// func main(){
// 	db.Init()
// 	handlers.InitCollection()
// 	r := gin.Default()

// 	// Apply CORS middleware
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = []string{"http://localhost:5173"} // Add your React app's origin
// 	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
// 	config.AllowHeaders = []string{"Authorization", "Content-Type"}
// 	r.Use(cors.New(config))

// 	handlers.SetupRoutes(r)
// 	r.Run(":8080")
// }