package models

import models "GoMJTrainingCamp/dbs/models/users"

type TrainingClassDetail struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	UserID          *uint         `gorm:"null" json:"user_id"`
	User            models.User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	TrainingClassID *uint         `gorm:"null;index" json:"training_class_id"` // Foreign key to TrainingClass
	TrainingClass   TrainingClass `gorm:"foreignKey:TrainingClassID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"training_class"`
}
