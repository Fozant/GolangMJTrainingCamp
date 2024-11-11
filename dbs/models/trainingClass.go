package models

import (
	"time"
)

type TrainingClass struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	ClassName        string    `gorm:"type:varchar(255);not null" json:"className"`
	ClassRequirement string    `gorm:"type:text" json:"classRequirement"`
	ClassNote        string    `gorm:"type:text" json:"classNote"`
	ClassDate        time.Time `gorm:"not null" json:"classDate"`
	ClassTime        time.Time `gorm:"not null" json:"classTime"`
	ClassCapacity    int64     `gorm:"not null" json:"classCapacity"`
}
