package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"errors"
)

func CreateTrainingClass(class *models.TrainingClass) error {
	if result := dbConnection.DB.Create(&class); result.Error != nil {
		return errors.New("database error")
	}
	return nil
}
