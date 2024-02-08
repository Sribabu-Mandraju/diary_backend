package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/admin/allAdmins", controllers.GetAllAdmins())
	incomingRoutes.GET("/admin/adminInfo", controllers.GetUserInfo())
	incomingRoutes.POST("/admin/sendRequest", controllers.SendRequest())
	incomingRoutes.PUT("/admin/modify-request/:id",controllers.ApproveOrRejectRequest())
	incomingRoutes.GET("/admin/all-requests",controllers.GetAllRequests())
	incomingRoutes.GET("/admin/adminsList", controllers.GetAllAdmins())
	incomingRoutes.GET("/admin/adminByID/:id", controllers.GetAdminByID())
	incomingRoutes.GET("/admin/client/all-clients", controllers.GetAllClients())
	incomingRoutes.GET("/admin/client/:id",controllers.GetClientByID())
}