package models

import (
	"fmt"
	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, userID uint) (*User, error) {
	var user User
	result := db.First(&user, userID) // Retrieves the first matching user based on userID
	if result.Error != nil {
		return nil, fmt.Errorf("user not found: %w", result.Error)
	}
	return &user, nil
}
