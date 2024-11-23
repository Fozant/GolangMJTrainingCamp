package models

import (
	"time"
)

type TrainingClass struct {
	IDTrainingClass  uint      `gorm:"primaryKey;autoIncrement;not null" json:"id_training_class"`
	ClassName        string    `gorm:"type:varchar(255);not null" json:"className"`
	ClassRequirement string    `gorm:"type:text" json:"classRequirement,omitempty"`
	ClassNote        string    `gorm:"type:text" json:"classNote,omitempty"`
	ClassDateTime    time.Time `gorm:"type:datetime;not null" json:"classDateTime"`
	ClassCapacity    uint      `gorm:"not null" json:"classCapacity"`
}
