package models

import (
	"time"
)

type TrainingClass struct {
	IDTrainingClass  uint      `gorm:"primaryKey;autoIncrement;type:uint unsigned" json:"IDTrainingClass"`
	ClassName        string    `gorm:"type:varchar(255);not null" json:"className"`
	ClassRequirement string    `gorm:"type:text" json:"classRequirement,omitempty"`
	ClassNote        string    `gorm:"type:text" json:"classNote,omitempty"`
	ClassDateTime    time.Time `gorm:"type:datetime;not null" json:"classDateTime"`
	ClassCapacity    int64     `gorm:"not null" json:"classCapacity"`
}
