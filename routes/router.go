package routes

import (
	authController "GoMJTrainingCamp/controller"
	membershipController "GoMJTrainingCamp/controller"
	packageListController "GoMJTrainingCamp/controller"
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	transactionController "GoMJTrainingCamp/controller"
	visitPackageController "GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, classHandler *trainingClassController.ClassHandler,
	trainerHandler *trainerController.TrainerHandler, membershipHandler *membershipController.MembershipHandler,
	visitHandler *visitPackageController.VisitHandler, transactionHandler *transactionController.TransactionHandler,
	packageHandler *packageListController.PackageHandler) {

	authRoutes(r)
	classRoutes(r, classHandler)
	trainerRoutes(r, trainerHandler)
	membershipRoutes(r, membershipHandler)
	visitRoutes(r, visitHandler)
	transactionRoutes(r, transactionHandler)
	packageRoutes(r, packageHandler)

}

func authRoutes(r *gin.Engine) {
	r.GET("/api/login", authController.HandleLogin)
	r.POST("/api/register", authController.HandleRegister)
}

func classRoutes(r *gin.Engine, handler *trainingClassController.ClassHandler) {
	// Protected class routes with JWT authentication
	classGroup := r.Group("/api/class")
	classGroup.Use(service.WithJWTAuth) // JWT middleware for these routes
	classGroup.POST("/add", handler.CreateClass)
	classGroup.GET("/get", handler.GetClasses)
	classGroup.GET("/getHistory", handler.GetClassesHistory)
	classGroup.POST("/book", handler.BookClass)
	classGroup.GET("/trainerschedule", handler.GetTrainerSchedule)
	classGroup.GET("/check-eligibility", handler.CheckEligibility)
}

func trainerRoutes(r *gin.Engine, handler *trainerController.TrainerHandler) {

	trainerGroup := r.Group("/api/trainer")
	trainerGroup.Use(service.WithJWTAuth)
	trainerGroup.POST("/add", handler.AddTrainer)
	trainerGroup.GET("/get", handler.GetAllTrainer)
}

func transactionRoutes(r *gin.Engine, handler *transactionController.TransactionHandler) {

	transactionGroup := r.Group("/api/transaction")
	transactionGroup.Use(service.WithJWTAuth)
	transactionGroup.POST("/verify", handler.VerifyTransaction)
	transactionGroup.GET("/get", handler.GetTransaction)
	transactionGroup.GET("/getbyuser", handler.GetTransactionByUser)
}

func membershipRoutes(r *gin.Engine, handler *membershipController.MembershipHandler) {
	membershipGroup := r.Group("/api/membership")
	membershipGroup.Use(service.WithJWTAuth)
	membershipGroup.POST("/buy", handler.BuyMembership)
}

func visitRoutes(r *gin.Engine, handler *visitPackageController.VisitHandler) {
	visitGroup := r.Group("/api/visit")
	visitGroup.Use(service.WithJWTAuth)
	visitGroup.POST("/buy", handler.BuyVisit)
}
func packageRoutes(r *gin.Engine, handler *packageListController.PackageHandler) {
	packageGroup := r.Group("/api/package")
	packageGroup.Use(service.WithJWTAuth)
	packageGroup.POST("/add", handler.AddPackage)
	packageGroup.GET("/get", handler.GetPackage)
}
