package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	users "GoMJTrainingCamp/dbs/models/users"
	"fmt"
)

type VisitServiceInterface interface {
	BuyVisit(visit *models.VisitPackage) (uint, error)
}
type VisitService struct {
}

func NewVisitService() VisitServiceInterface {
	return &VisitService{}
}

func (s *VisitService) BuyVisit(visit *models.VisitPackage) (uint, error) {
	var user users.User

	if err := dbConnection.DB.First(&user, visit.UserID).Error; err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	if err := dbConnection.DB.Create(&visit).Error; err != nil {
		return 0, fmt.Errorf("failed to create membership: %w", err)
	}
	return visit.IDVisitPackage, nil
}
