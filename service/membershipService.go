package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	users "GoMJTrainingCamp/dbs/models/users"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MembershipWithTransaction struct {
	MembershipID  uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	TransactionID uint      `json:"transaction_id"`
	PaymentStatus string    `json:"paymentStatus"`
}

type MembershipServiceInterface interface {
	BuyMembership(membership *models.Membership) (uint, error)
	UpdateTransactionID(membershipID uint, transactionID uint) error
	GetMembershipByUser(userID uint) ([]MembershipWithTransaction, error)
}
type MembershipService struct {
}

func NewMembershipService() MembershipServiceInterface {
	return &MembershipService{}
}

func (s *MembershipService) BuyMembership(membership *models.Membership) (uint, error) {
	var user users.User

	if err := dbConnection.DB.First(&user, membership.UserID).Error; err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	if err := dbConnection.DB.Create(&membership).Error; err != nil {
		return 0, fmt.Errorf("failed to create membership: %w", err)
	}
	return membership.IDMembership, nil
}
func (s *MembershipService) UpdateTransactionID(membershipID uint, transactionID uint) error {
	return dbConnection.DB.Model(&models.Membership{}).
		Where("id_membership = ?", membershipID).
		Update("id_transaction", transactionID).Error
}
func (s *MembershipService) GetMembershipByUser(userID uint) ([]MembershipWithTransaction, error) {
	var results []MembershipWithTransaction

	err := dbConnection.DB.Raw(`
    SELECT memberships.*, transactions.*
    FROM memberships
    LEFT JOIN transactions ON memberships.id_membership = transactions.membership_id
    WHERE memberships.user_id = ?
`, userID).Scan(&results).Error
	fmt.Println(results)
	// Check for specific database errors
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Handle case when no records are found for the user
			return nil, fmt.Errorf("no memberships found for user with ID %d: %w", userID, err)
		}

		// Handle other types of errors such as database connection issues
		return nil, fmt.Errorf("error querying memberships for user with ID %d: %w", userID, err)
	}

	// If no results are found, handle gracefully
	if len(results) == 0 {
		return nil, fmt.Errorf("no memberships found for user with ID %d", userID)
	}

	// Return the results
	return results, nil
}
