package routes

import (
	"backend/controllers"
	middleware "backend/middlewares"

	"github.com/gin-gonic/gin"
)

func DiaryRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("diary/getAllDiaries/:userid", controllers.GetAllDiaries())
	incomingRoutes.GET("diary/getDiary/:id", controllers.GetDiaryById())
	incomingRoutes.POST("diary/createDiary", controllers.CreateDiary())
	incomingRoutes.PATCH("diary/updateDiary/:id", controllers.UpdateDiary())
	incomingRoutes.DELETE("diary/deleteDiary/:id", controllers.DeleteDiary())
}
