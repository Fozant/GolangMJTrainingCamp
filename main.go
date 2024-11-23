package main

import (
	"GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/dbs/models/trainer"
	users "GoMJTrainingCamp/dbs/models/users"
	"GoMJTrainingCamp/routes"
	"GoMJTrainingCamp/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	// Display a custom banner at startup
	displayBanner()

	// Attempt to connect to the database and display status
	startTime := time.Now()
	if err := dbConnection.ConnectDatabase(); err != nil {
		log.Fatalf("❌ Database Connection Failed: %v", err)
	} else {
		fmt.Println("✅ Database Connection Successful!")
	}

	if err := dbConnection.DB.AutoMigrate(&users.User{}, &models.TrainingClass{}, trainer.Trainer{}, &models.Transaction{}, &models.Membership{}, &models.VisitPackage{}); err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		return
	}

	if err := dbConnection.DB.AutoMigrate(&models.TrainingClassDetail{}); err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		return
	}

	classHandler, trainerHandler, membershipHandler, visitHandler, transactionHandler := initHandler()
	r := gin.Default()
	routes.SetupRoutes(r, classHandler, trainerHandler, membershipHandler, visitHandler, transactionHandler)

	// Run the server and display server details
	port := ":8080"

	fmt.Printf("🚀 Starting application on http://localhost%s\n", port)
	fmt.Printf("Application started in %v seconds\n", time.Since(startTime).Seconds())
	if err := r.Run(port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}

}

func displayBanner() {
	fmt.Println("===========================================")
	fmt.Println("     MY GO APP - REST API WITH GIN & GORM  ")
	fmt.Println("===========================================")
	fmt.Println(":: My Go Application ::       (v0.1-SNAPSHOT)")
	fmt.Println("===========================================")
	fmt.Println("Starting Application...")
	fmt.Println()
}

func initHandler() (
	*controller.ClassHandler,
	*controller.TrainerHandler,
	*controller.MembershipHandler,
	*controller.VisitHandler,
	*controller.TransactionHandler,
) {
	// Initialize the services
	classService := service.NewClassService()
	trainerService := service.NewTrainerService()
	membershipService := service.NewMembershipService()
	transactionService := service.NewTransactionService()
	visitPackageService := service.NewVisitService()

	// Initialize the handlers with the corresponding services
	classHandler := controller.NewClassHandler(classService, membershipService, visitPackageService)
	trainerHandler := controller.NewTrainerHandler(trainerService)
	membershipHandler := controller.NewMembershipHandler(membershipService, transactionService)
	visitHandler := controller.NewVisitHandler(visitPackageService, transactionService)
	transactionHandler := controller.NewTransactionHandler(transactionService)

	// Return the handlers to be used in the main function
	return classHandler, trainerHandler, membershipHandler, visitHandler, transactionHandler
}
