package routes

import (
	authController "GoMJTrainingCamp/controller"
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, classHandler *trainingClassController.
	ClassHandler, trainerHandler *trainerController.TrainerHandler) {
	authRoutes(r)
	classRoutes(r, classHandler)
	trainerRoutes(r, trainerHandler)

}

func authRoutes(r *gin.Engine) {
	r.GET("/api/login", authController.HandleLogin)
	r.POST("/api/register", authController.HandleRegister)
}

func classRoutes(r *gin.Engine, handler *trainingClassController.ClassHandler) {
	classGroup := r.Group("/api/class")
	classGroup.Use(service.WithJWTAuth)
	classGroup.POST("/add", handler.CreateClass)
	classGroup.GET("/get", handler.GetClasses)
}
func trainerRoutes(r *gin.Engine, handler *trainerController.TrainerHandler) {
	trainerGroup := r.Group("/api/trainer")
	trainerGroup.Use(service.WithJWTAuth)
	trainerGroup.POST("/add", handler.AddTrainer)
}
