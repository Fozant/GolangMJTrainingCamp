package routes

import (
	authController "GoMJTrainingCamp/controller"
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authRoutes(r)
	classRoutes(r)
	trainerRoutes(r)
}

func authRoutes(r *gin.Engine) {
	r.GET("/api/login", authController.HandleLogin)
	r.POST("/api/register", authController.HandleRegister)
}

func classRoutes(r *gin.Engine) {

	classGroup := r.Group("/api/class")
	classGroup.Use(service.WithJWTAuth)
	classGroup.POST("/add", trainingClassController.CreateClass)
	classGroup.GET("/get", trainingClassController.GetClasses)
}
func trainerRoutes(r *gin.Engine) {

	classGroup := r.Group("/api/trainer")
	classGroup.Use(service.WithJWTAuth)
	classGroup.POST("/add", trainerController.AddTrainer)
}
