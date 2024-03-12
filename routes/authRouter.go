package routes

import (
	"backend/controllers" // Importing controllers package for handling HTTP request handlers

	"github.com/gin-gonic/gin" // Importing Gin web framework package
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// AuthRoutes function defines authentication-related routes and registers them with the provided Gin engine.

	incomingRoutes.POST("admin/login", controllers.Login())
	// Defines a route to handle user login requests.

	incomingRoutes.POST("admin/signup", controllers.Register())
	// Defines a route to handle user signup/registration requests.
}
