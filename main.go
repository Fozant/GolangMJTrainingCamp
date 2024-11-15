package main

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	trainer "GoMJTrainingCamp/dbs/models/trainer"
	users "GoMJTrainingCamp/dbs/models/users"
	"GoMJTrainingCamp/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	// Display a custom banner at startup
	displayBanner()

	// Initialize Gin router
	r := gin.Default()

	// Attempt to connect to the database and display status
	startTime := time.Now()
	if err := dbConnection.ConnectDatabase(); err != nil {
		log.Fatalf("❌ Database Connection Failed: %v", err)
	} else {
		fmt.Println("✅ Database Connection Successful!")
	}

	if err := dbConnection.DB.AutoMigrate(&users.User{}, &models.TrainingClass{}, trainer.Trainer{}); err != nil {
		fmt.Printf("failed to migrate database: %v\n", err)
		return // Exiting the function after logging the error
	}

	// Set up routes via the routes package
	routes.SetupRoutes(r)

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
