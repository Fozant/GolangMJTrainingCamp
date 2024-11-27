package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	models "GoMJTrainingCamp/dbs/models/users"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ValidateUserCredentials validates user credentials (email and password)
func ValidateUserCredentials(email, password string) (*models.User, error) {
	var user models.User

	// Query the user by email
	if err := dbConnection.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Query the user by email
	if err := dbConnection.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	// Create the user in the database
	if err := dbConnection.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}
