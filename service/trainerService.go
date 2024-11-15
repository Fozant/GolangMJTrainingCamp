package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	trainer2 "GoMJTrainingCamp/dbs/models/trainer"
	"errors"
	"log"
)

func AddTrainer(trainer *trainer2.Trainer) (idTrainer uint, err error) {
	if result := dbConnection.DB.Create(&trainer); result.Error != nil {
		log.Printf("Error creating class: %v", result.Error)
		return 0, errors.New("database error")
	}
	return trainer.IDTrainer, nil
}
