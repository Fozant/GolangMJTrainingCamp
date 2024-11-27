package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
)

type PackageServiceInterface interface {
	CreatePackage(packages *models.PackageList) error
	GetPackage(id string) ([]models.PackageList, error)
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
func (service *PackageService) GetPackage(id string) ([]models.PackageList, error) {
	var packages []models.PackageList

	if id != "" {
		err := dbConnection.DB.Where("id_package = ?", id).Find(&packages).Error
		if err != nil {
			return nil, err
		}
		return packages, nil
	}

	err := dbConnection.DB.Find(&packages).Error
	if err != nil {
		return nil, err
	}
	return packages, nil
}
