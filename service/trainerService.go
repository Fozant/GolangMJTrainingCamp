package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	trainer2 "GoMJTrainingCamp/dbs/models/trainer"
	"errors"
	"log"
)

type TrainerServiceInterface interface {
	AddTrainer(trainer *trainer2.Trainer) (idTrainer uint, err error)
	GetTrainer() ([]trainer2.Trainer, error)
}
type TrainerService struct {
}

func NewTrainerService() TrainerServiceInterface {
	return &TrainerService{}
}

func (s *TrainerService) AddTrainer(trainer *trainer2.Trainer) (idTrainer uint, err error) {
	if result := dbConnection.DB.Create(&trainer); result.Error != nil {
		log.Printf("Error creating trainer: %v", result.Error)
		return 0, errors.New("database error")
	}

	return trainer.ID, nil
}

func (s *TrainerService) GetTrainer() ([]trainer2.Trainer, error) {
	var trainers []trainer2.Trainer

	result := dbConnection.DB.Find(&trainers)
	if result.Error != nil {
		log.Printf("Error getting trainers: %v", result.Error)
		return nil, result.Error
	}

	return trainers, nil
}
