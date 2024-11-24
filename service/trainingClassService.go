package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/dbs/models/trainer"
	users "GoMJTrainingCamp/dbs/models/users"
	"errors"
	"fmt"
	"log"
	"time"
)

type GetClassResponse struct {
	IDClass          uint          `json:"idClass"`
	ClassDate        time.Time     `json:"classDate"`
	ClassCapacity    uint          `json:"classCapacity"`
	ClassRequirement string        `json:"classRequirement"`
	ClassName        string        `json:"className"`
	TrainerDetail    TrainerDetail `json:"trainerDetail"`
	ClassMembers     []ClassMember `json:"classMembers"`
}

type TrainerDetail struct {
	IDTrainer          uint   `json:"idTrainer"`
	TrainerName        string `json:"trainerName"`
	TrainerDescription string `json:"trainerDescription"`
}

type ClassMember struct {
	IDUser  uint   `json:"idUser"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	PNumber string `json:"pNumber"`
}

type ClassServiceInterface interface {
	CreateTrainingClass(class *models.TrainingClass) error
	GetClasses(id, date string) ([]GetClassResponse, error)
	GetClassesHistory(userID uint) ([]GetClassResponse, error)
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

func (s *ClassService) GetClasses(id, date string) ([]GetClassResponse, error) {
	var responses []GetClassResponse
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

	for _, class := range classes {
		var trainer trainer.Trainer
		if err := dbConnection.DB.Where("id =?", class.TrainerID).First(&trainer).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve trainer for class ID %d: %w", class.IDTrainingClass, err)
		}
		var members []models.TrainingClassDetail
		if err := dbConnection.DB.Where("training_class_id=?", class.IDTrainingClass).Find(&members).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve members for class ID %d: %w", class.IDTrainingClass, err)
		}
		var classMembers []ClassMember
		if len(members) > 0 {
			for _, member := range members {
				var user users.User
				if err := dbConnection.DB.Where("id_user=?", member.UserID).First(&user).Error; err != nil {
					return nil, fmt.Errorf("failed to retrieve user details for member ID %d: %w", member.UserID, err)
				}
				classMembers = append(classMembers, ClassMember{
					IDUser:  user.IDUser,
					Name:    user.Name,
					Email:   user.Email,
					PNumber: user.PNumber,
				})
			}
		} else {
			// If no members, initialize with an empty slice (optional but good practice)
			classMembers = []ClassMember{}
		}
		responses = append(responses, GetClassResponse{
			IDClass:          class.IDTrainingClass,
			ClassDate:        class.ClassDateTime,
			ClassCapacity:    class.ClassCapacity,
			ClassRequirement: class.ClassRequirement,
			ClassName:        class.ClassName,
			TrainerDetail: TrainerDetail{
				IDTrainer:          trainer.ID,
				TrainerName:        trainer.TrainerName,
				TrainerDescription: trainer.TrainerDescription},
			ClassMembers: classMembers,
		})
	}

	return responses, nil
}
func (s *ClassService) GetClassesHistory(userID uint) ([]GetClassResponse, error) {
	var responses []GetClassResponse
	var classDetails []models.TrainingClassDetail

	// Find class details where the user has participated
	if err := dbConnection.DB.Where("user_id = ?", userID).Find(&classDetails).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve class history for user ID %d: %w", userID, err)
	}
	if len(classDetails) == 0 {
		return nil, fmt.Errorf("no class history found for user ID %d", userID)
	}

	// Retrieve unique class IDs from the user's class details
	classIDMap := make(map[uint]bool)
	for _, detail := range classDetails {
		classIDMap[*detail.TrainingClassID] = true
	}

	// Retrieve the actual class information
	var classIDs []uint
	for id := range classIDMap {
		classIDs = append(classIDs, id)
	}

	var classes []models.TrainingClass
	if err := dbConnection.DB.Where("id_training_class IN (?)", classIDs).Find(&classes).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve class details: %w", err)
	}

	for _, class := range classes {
		// Retrieve trainer information
		var trainer trainer.Trainer
		if err := dbConnection.DB.Where("id = ?", class.TrainerID).First(&trainer).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve trainer for class ID %d: %w", class.IDTrainingClass, err)
		}

		// Retrieve all members for this class
		var members []models.TrainingClassDetail
		if err := dbConnection.DB.Where("training_class_id = ?", class.IDTrainingClass).Find(&members).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve members for class ID %d: %w", class.IDTrainingClass, err)
		}

		// Populate class members
		var classMembers []ClassMember
		for _, member := range members {
			var user users.User
			if err := dbConnection.DB.Where("id_user = ?", member.UserID).First(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to retrieve user details for member ID %d: %w", member.UserID, err)
			}
			classMembers = append(classMembers, ClassMember{
				IDUser:  user.IDUser,
				Name:    user.Name,
				Email:   user.Email,
				PNumber: user.PNumber,
			})
		}

		// Append the class to the response
		responses = append(responses, GetClassResponse{
			IDClass:          class.IDTrainingClass,
			ClassDate:        class.ClassDateTime,
			ClassCapacity:    class.ClassCapacity,
			ClassRequirement: class.ClassRequirement,
			ClassName:        class.ClassName,
			TrainerDetail: TrainerDetail{
				IDTrainer:          trainer.ID,
				TrainerName:        trainer.TrainerName,
				TrainerDescription: trainer.TrainerDescription,
			},
			ClassMembers: classMembers,
		})
	}

	return responses, nil
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
