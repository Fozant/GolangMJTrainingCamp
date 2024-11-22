package routes

import (
	authController "GoMJTrainingCamp/controller"
	membershipController "GoMJTrainingCamp/controller"
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	visitPackageController "GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, classHandler *trainingClassController.
	ClassHandler, trainerHandler *trainerController.TrainerHandler, membershipHandler *membershipController.MembershipHandler,
	visitHandler *visitPackageController.VisitHandler) {
	authRoutes(r)
	classRoutes(r, classHandler)
	trainerRoutes(r, trainerHandler)
	membershipRoutes(r, membershipHandler)
	visitRoutes(r, visitHandler)

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
func membershipRoutes(r *gin.Engine, handler *membershipController.MembershipHandler) {
	trainerGroup := r.Group("/api/membership")
	trainerGroup.Use(service.WithJWTAuth)
	trainerGroup.POST("/buy", handler.BuyMembership)
}
func visitRoutes(r *gin.Engine, handler *visitPackageController.VisitHandler) {
	visitGroup := r.Group("/api/visit")
	visitGroup.Use(service.WithJWTAuth)
	visitGroup.POST("/buy", handler.BuyVisit)

}
