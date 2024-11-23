package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"errors"
	"fmt"
	"log"
)

type ClassServiceInterface interface {
	CreateTrainingClass(class *models.TrainingClass) error
	GetClasses(id, date string) ([]models.TrainingClass, error)
	BookClass(classDetail *models.TrainingClassDetail) error
	AlreadyBooked(userID, trainingClassID uint) (bool, error)
	CountParticipant(trainingClassID uint) (uint, error)
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
		if err := dbConnection.DB.Where("id_training_class = ?", id).Find(&classes).Error; err != nil {
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
	if len(classes) == 0 {
		return nil, fmt.Errorf("no classes found ")
	}
	return classes, nil
}
func (s *ClassService) BookClass(classDetail *models.TrainingClassDetail) error {

	if result := dbConnection.DB.Create(&classDetail); result.Error != nil {
		log.Printf("Error  : %v", result.Error)
		return errors.New("database error")
	}
	return nil
}

func (s *ClassService) AlreadyBooked(userID, trainingClassID uint) (bool, error) {
	var exists bool
	err := dbConnection.DB.Raw(
		"SELECT EXISTS (SELECT 1 FROM training_class_details WHERE user_id = ? AND training_class_id = ?)",
		userID, trainingClassID,
	).Scan(&exists).Error

	if err != nil {
		log.Printf("Database error: %v", err)
		return false, err
	}
	return exists, nil
}

func (s *ClassService) CountParticipant(trainingClassID uint) (uint, error) {
	var total uint
	err := dbConnection.DB.Raw(
		"SELECT COUNT(*) FROM training_class_details WHERE training_class_id = ?",
		trainingClassID,
	).Scan(&total).Error

	if err != nil {
		log.Printf("Database error: %v", err)
		return 0, err
	}
	return total, nil
}
