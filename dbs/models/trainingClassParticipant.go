package models

import models "GoMJTrainingCamp/dbs/models/users"

type TrainingClassDetail struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	UserID          *uint         `gorm:"null" json:"user_id"`
	User            models.User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	IDTrainingClass *uint         `gorm:"null" json:"IDTrainingClass"`
	TrainingClass   TrainingClass `gorm:"foreignKey:IDTrainingClass;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"training_class"`
}
