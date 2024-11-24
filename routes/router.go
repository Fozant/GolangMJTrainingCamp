package routes

import (
	authController "GoMJTrainingCamp/controller"
	membershipController "GoMJTrainingCamp/controller"
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	transactionController "GoMJTrainingCamp/controller"
	visitPackageController "GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, classHandler *trainingClassController.ClassHandler,
	trainerHandler *trainerController.TrainerHandler, membershipHandler *membershipController.MembershipHandler,
	visitHandler *visitPackageController.VisitHandler, transactionHandler *transactionController.TransactionHandler) {

	// Public routes for login and registration
	authRoutes(r)

	// Protected routes with JWT authentication
	classRoutes(r, classHandler)
	trainerRoutes(r, trainerHandler)
	membershipRoutes(r, membershipHandler)
	visitRoutes(r, visitHandler)
	transactionRoutes(r, transactionHandler)
}

func authRoutes(r *gin.Engine) {
	// Public authentication routes
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
}

func trainerRoutes(r *gin.Engine, handler *trainerController.TrainerHandler) {
	// Protected trainer routes with JWT authentication
	trainerGroup := r.Group("/api/trainer")
	trainerGroup.Use(service.WithJWTAuth) // JWT middleware for these routes
	trainerGroup.POST("/add", handler.AddTrainer)
}

func transactionRoutes(r *gin.Engine, handler *transactionController.TransactionHandler) {
	// Protected transaction routes with JWT authentication
	transactionGroup := r.Group("/api/transaction")
	transactionGroup.Use(service.WithJWTAuth) // JWT middleware for these routes
	transactionGroup.POST("/verify", handler.VerifyTransaction)
	transactionGroup.GET("/get", handler.GetTransaction)
}

func membershipRoutes(r *gin.Engine, handler *membershipController.MembershipHandler) {
	// Protected membership routes with JWT authentication
	membershipGroup := r.Group("/api/membership")
	membershipGroup.Use(service.WithJWTAuth) // JWT middleware for these routes
	membershipGroup.POST("/buy", handler.BuyMembership)
}

func visitRoutes(r *gin.Engine, handler *visitPackageController.VisitHandler) {
	// Protected visit routes with JWT authentication
	visitGroup := r.Group("/api/visit")
	visitGroup.Use(service.WithJWTAuth) // JWT middleware for these routes
	visitGroup.POST("/buy", handler.BuyVisit)
}
