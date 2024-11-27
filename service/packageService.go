package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
)

type PackageServiceInterface interface {
	//BuyMembership(membership *models.Membership) (uint, error)
	//UpdateTransactionID(membershipID uint, transactionID uint) error
	//GetMembershipByUser(userID uint) ([]MembershipWithTransaction, error)
	CreatePackage(packages *models.PackageList) error
}
type PackageService struct {
}

func NewPackageService() PackageServiceInterface {
	return &PackageService{}
}

func (service *PackageService) CreatePackage(packages *models.PackageList) error {

	result := dbConnection.DB.Create(&packages)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
