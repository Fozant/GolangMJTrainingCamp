package models

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models/trainer"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Role string

const (
	RoleUser    Role = "ROLE_USER"
	RoleAdmin   Role = "ROLE_ADMIN"
	RoleTrainer Role = "ROLE_TRAINER"
)

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if !user.isValidRole() {
		return errors.New("invalid role")
	}
	return nil
}

func (user *User) isValidRole() bool {
	switch user.Role {
	case RoleUser, RoleAdmin, RoleTrainer:
		return true
	}
	return false
}

// User model
type User struct {
	IDUser           uint             `gorm:"primaryKey;autoIncrement" json:"id_user"`
	PNumber          string           `gorm:"type:varchar(50);not null" json:"p_number"`
	Name             string           `gorm:"type:varchar(100);not null" json:"name"`
	Email            string           `gorm:"type:varchar(100);unique;not null" json:"email"`
	RegistrationDate time.Time        `gorm:"type:datetime;not null" json:"registration_date"`
	Password         string           `gorm:"type:varchar(255);not null" json:"password"`
	Role             Role             `gorm:"type:varchar(50);not null" json:"role"`
	IDTrainer        *trainer.Trainer `gorm:"foreignKey:IDTrainer;constraint:OnDelete:SET NULL;" json:"id_trainer,omitempty"`
}

func (u *User) GetAuthorities() []string {
	return []string{string(u.Role)}
}

func GetUserByID(userID uint) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE id_user = ?"

	result := dbConnection.DB.Raw(query, userID).Scan(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found: %w", result.Error)
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	// Query the database to find the user by email
	if err := dbConnection.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user User) error {

	if err := dbConnection.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}
