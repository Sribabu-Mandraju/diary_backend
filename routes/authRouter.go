package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("admin/login", controllers.Login())
	incomingRoutes.POST("admin/signup", controllers.Register())
	
}

func ClientAuthRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("client/login", controllers.ClientLogin())
	incomingRoutes.POST("client/signup", controllers.ClientRegister())
}



