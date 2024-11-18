package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"errors"
	"log"
)

type ClassServiceInterface interface {
	CreateTrainingClass(class *models.TrainingClass) error
	GetClasses(id, date string) ([]models.TrainingClass, error)
}
type ClassService struct{}

func NewClassService() ClassServiceInterface {
	return &ClassService{}
}

func (s *ClassService) CreateTrainingClass(class *models.TrainingClass) error {
	if result := dbConnection.DB.Create(&class); result.Error != nil {
		log.Printf("Error creating class: %v", result.Error)
		return errors.New("database error")
	}
	return nil
}

func (s *ClassService) GetClasses(id, date string) ([]models.TrainingClass, error) {
	var classes []models.TrainingClass

	if id != "" {
		if err := dbConnection.DB.Where("id = ?", id).Find(&classes).Error; err != nil {
			return nil, err
		}
	} else if date != "" {
		if err := dbConnection.DB.Where("DATE(Class_date_time) = ?", date).Find(&classes).Error; err != nil {
			return nil, err
		}
	} else {
		if err := dbConnection.DB.Find(&classes).Error; err != nil {
			return nil, err
		}
	}
	return classes, nil
}
