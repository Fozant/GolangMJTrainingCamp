package main

import (
	membershipController "GoMJTrainingCamp/controller"
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	transactionController "GoMJTrainingCamp/controller"
	visitPackageController "GoMJTrainingCamp/controller"

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

	if err := dbConnection.DB.AutoMigrate(&users.User{}, &models.TrainingClass{}, trainer.Trainer{}, &models.Membership{}, &models.Transaction{}); err != nil {
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
	*trainingClassController.ClassHandler,
	*trainerController.TrainerHandler,
	*membershipController.MembershipHandler,
	*visitPackageController.VisitHandler,
	*transactionController.TransactionHandler,
) {

	classService := service.NewClassService()
	trainerService := service.NewTrainerService()
	membershipService := service.NewMembershipService()
	transactionService := service.NewTransactionService()
	visitPackageService := service.NewVisitService()

	classHandler := trainingClassController.NewClassHandler(classService)
	trainerHandler := trainerController.NewTrainerHandler(trainerService)
	membershipHandler := membershipController.NewMembershipHandler(membershipService, transactionService)
	visitHandler := visitPackageController.NewVisitHandler(visitPackageService, transactionService)
	transactionHandler := transactionController.NewTransactionHandler(transactionService)

	return classHandler, trainerHandler, membershipHandler, visitHandler, transactionHandler
}

func automigrate() string {

	return "&users.User{}, &models.TrainingClass{}, trainer.Trainer{}, &models.Membership{},&models.Transaction{}"
}
