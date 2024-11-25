package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	users "GoMJTrainingCamp/dbs/models/users"
	"fmt"
	"gorm.io/gorm"
)

type VisitWithTransaction struct {
	VisitID       uint   `json:"id" gorm:"column:id_visit_package"`
	UserID        uint   `json:"user_id" gorm:"column:user_id"`
	VisitNumber   uint   `json:"visitNumber" gorm:"column:visit_number"`
	VisitUsed     uint   `json:"visitUsed" gorm:"column:visit_used"`
	TransactionID uint   `json:"transaction_id" gorm:"column:id_transaction"`
	PaymentStatus string `json:"paymentStatus" gorm:"column:paymentStatus"`
}
type VisitServiceInterface interface {
	BuyVisit(visit *models.VisitPackage) (uint, error)
	GetVisitByUser(userID uint) ([]VisitWithTransaction, error)
	UseVisit(userID uint, idVisit uint) error
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

func (s *VisitService) GetVisitByUser(userID uint) ([]VisitWithTransaction, error) {
	var results []VisitWithTransaction

	err := dbConnection.DB.Raw(`
    SELECT visit_packages.*, transactions.*
    FROM visit_packages
    LEFT JOIN transactions ON visit_packages.id_visit_package = transactions.Visit_id
    WHERE visit_packages.user_id = ?
`, userID).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Handle case when no records are found for the user
			return nil, fmt.Errorf("no visit found for user with ID %d: %w", userID, err)
		}

		// Handle other types of errors such as database connection issues
		return nil, fmt.Errorf("error querying visit for user with ID %d: %w", userID, err)
	}

	// If no results are found, handle gracefully
	if len(results) == 0 {
		return nil, fmt.Errorf("no visit found for user with ID %d", userID)
	}

	// Return the results
	return results, nil
}

func (s *VisitService) UseVisit(userID uint, idVisit uint) error {

	var visit models.VisitPackage
	if err := dbConnection.DB.Model(&models.VisitPackage{}).Where("id_visit_package = ?",
		idVisit).First(&visit).Error; err != nil {
		return fmt.Errorf("failed to find visit package: %w", err)
	}
	if err := dbConnection.DB.Model(&models.VisitPackage{}).Where("id_visit_package = ?", idVisit).Updates(map[string]interface{}{
		"visit_used": visit.VisitUsed + 1, // Increment the `VisitUsed` field by 1
	}).Error; err != nil {
		return fmt.Errorf("failed to update visit package: %w", err)
	}
	return nil
}
