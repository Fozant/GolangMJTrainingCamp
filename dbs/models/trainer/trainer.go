package trainer

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"gorm.io/gorm"
	"log"
)

type Trainer struct {
	gorm.Model
	TrainerName        string `gorm:"type:varchar(100);not null" json:"trainer_name"`
	TrainerDescription string `gorm:"type:text" json:"trainer_description"`
}

func CreateTrainer(trainer Trainer) error {

	if err := dbConnection.DB.Create(&trainer).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}
