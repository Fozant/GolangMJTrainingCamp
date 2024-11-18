package main

import (
	trainerController "GoMJTrainingCamp/controller"
	trainingClassController "GoMJTrainingCamp/controller"
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	trainer "GoMJTrainingCamp/dbs/models/trainer"
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

	if err := dbConnection.DB.AutoMigrate(&users.User{}, &models.TrainingClass{}, trainer.Trainer{}); err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		return
	}
	classHandler, trainerHandler := initHandler()
	r := gin.Default()
	routes.SetupRoutes(r, classHandler, trainerHandler)

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
) {

	classService := service.NewClassService()
	trainerService := service.NewTrainerService()

	classHandler := trainingClassController.NewClassHandler(classService)
	trainerHandler := trainerController.NewTrainerHandler(trainerService)

	return classHandler, trainerHandler
}
