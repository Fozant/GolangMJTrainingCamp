package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	users "GoMJTrainingCamp/dbs/models/users"
	"fmt"
)

type MembershipServiceInterface interface {
	BuyMembership(membership *models.Membership) error
}
type MembershipService struct {
}

func NewMembershipService() MembershipServiceInterface {
	return &MembershipService{}
}

func (s *MembershipService) BuyMembership(membership *models.Membership) error {
	var user users.User

	if err := dbConnection.DB.First(&user, membership.UserID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := dbConnection.DB.Create(&membership).Error; err != nil {
		return fmt.Errorf("failed to create membership: %w", err)
	}
	return nil
}
