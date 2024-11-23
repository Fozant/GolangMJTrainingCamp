package models

import (
	models "GoMJTrainingCamp/dbs/models/users"
)

type VisitPackage struct {
	IDVisitPackage uint        `gorm:"primaryKey;autoIncrement" json:"IDVisitPackage"`
	UserID         uint        `json:"user_id"`
	User           models.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Price          uint        `json:"price" gorm:"type:int unsigned;not null"`
	VisitNumber    uint        `json:"visitNumber" gorm:"type:tinyint"`
	VisitUsed      uint        `json:"visitUsed" gorm:"type:tinyint"`
}
