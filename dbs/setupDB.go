package dbs

import (
	"GoMJTrainingCamp/dbs/model/users"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to the database and provides a custom error message
func ConnectDatabase() error {
	dsn := "root:@tcp(localhost:3306)/golangMJTC"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %w", err)
	}

	// Auto-migrate the schema
	if err := database.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Set the global DB variable
	DB = database
	return nil
}
