package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("admin/login", controllers.Login())
	incomingRoutes.POST("admin/signup", controllers.Register())
	incomingRoutes.GET("admin/info",controllers.GetUserInfo())
}




