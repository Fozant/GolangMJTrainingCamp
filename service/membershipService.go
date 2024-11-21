package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	users "GoMJTrainingCamp/dbs/models/users"
	"fmt"
)

type MembershipServiceInterface interface {
	BuyMembership(membership *models.Membership) (uint, error)
	UpdateTransactionID(membershipID uint, transactionID uint) error
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
