package routes

import (
	"backend/controllers"            // Importing controllers package for handling HTTP request handlers
	middleware "backend/middlewares" // Importing middlewares package for authentication middleware

	"github.com/gin-gonic/gin" // Importing Gin web framework package
)

func DiaryRoutes(incomingRoutes *gin.Engine) {
	// DiaryRoutes function defines routes related to diary operations and registers them with the provided Gin engine.

	incomingRoutes.Use(middleware.Authenticate()) // Using authentication middleware for all diary routes
	// This ensures that authentication is required for accessing the diary routes.

	incomingRoutes.GET("diary/getAllDiaries/:userid", controllers.GetAllDiaries())
	// Defines a route to retrieve all diaries belonging to a specific user by user ID.

	incomingRoutes.GET("diary/getDiary/:id", controllers.GetDiaryById())
	// Defines a route to retrieve a single diary by its ID.

	incomingRoutes.POST("diary/createDiary", controllers.CreateDiary())
	// Defines a route to create a new diary.

	incomingRoutes.PATCH("diary/updateDiary/:id", controllers.UpdateDiary())
	// Defines a route to update an existing diary by its ID.

	incomingRoutes.DELETE("diary/deleteDiary/:id", controllers.DeleteDiary())
	// Defines a route to delete a diary by its ID.
}
