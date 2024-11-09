package models

import (
	"GoMJTrainingCamp/dbs/model/trainer"
	"time"
)

type Role string

const (
	RoleUser    Role = "ROLE_USER"
	RoleAdmin   Role = "ROLE_ADMIN"
	RoleTrainer Role = "ROLE_TRAINER"
)

type User struct {
	IDUser           uint             `gorm:"primaryKey;autoIncrement" json:"id_user"`
	PNumber          string           `gorm:"type:varchar(50);not null" json:"p_number"`
	Name             string           `gorm:"type:varchar(100);not null" json:"name"`
	Email            string           `gorm:"type:varchar(100);unique;not null" json:"email"`
	RegistrationDate time.Time        `gorm:"not null" json:"registration_date"`
	Password         string           `gorm:"type:varchar(255);not null" json:"password"`
	Role             Role             `gorm:"type:varchar(50);not null" json:"role"`
	IDTrainer        *trainer.Trainer `gorm:"foreignKey:IDTrainer;constraint:OnDelete:SET NULL;" json:"id_trainer,omitempty"`
}

func (u *User) GetAuthorities() []string {
	return []string{string(u.Role)}
}
