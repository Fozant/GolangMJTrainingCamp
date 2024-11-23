package models

import (
	models "GoMJTrainingCamp/dbs/models/users"
	"time"
)

type Membership struct {
	IDMembership uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint        `json:"user_id"`
	User         models.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	StartDate    time.Time   `gorm:"type:datetime;not null" json:"startDate"`
	EndDate      time.Time   `gorm:"type:datetime;not null" json:"endDate"`
	Price        uint        `json:"price" gorm:"type:int unsigned;not null"`
	Duration     uint        `json:"duration" gorm:"type:tinyint"`
}
