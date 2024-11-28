package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CreateClassRequest struct {
	ClassName        string    `json:"className" binding:"required"`
	ClassRequirement string    `json:"classRequirement,omitempty"`
	ClassDateTime    time.Time `json:"ClassDateTime" binding:"required"`
	ClassCapacity    uint      `json:"classCapacity" binding:"required"`
	TrainerID        uint      `json:"trainerID" binding:"required"`
}

type BookClassRequest struct {
	IDClass uint   `json:"idClass" binding:"required"`
	IDUser  uint   `json:"idUser,omitempty"`
	Type    string `json:"type" binding:"required"`
}
type EligibilityResponse struct {
	ValidMember   bool `json:"validMember"`
	ValidVisit    bool `json:"validVisit"`
	AlreadyBooked bool `json:"alreadyBooked"`
}

type ClassHandler struct {
	ClassService      service.ClassServiceInterface
	MembershipService service.MembershipServiceInterface
	VisitService      service.VisitServiceInterface
}

func NewClassHandler(classService service.ClassServiceInterface, membershipService service.MembershipServiceInterface, VisitService service.VisitServiceInterface) *ClassHandler {
	return &ClassHandler{ClassService: classService,
		MembershipService: membershipService,
		VisitService:      VisitService}
}

func (h *ClassHandler) CreateClass(c *gin.Context) {
	var request CreateClassRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	class := models.TrainingClass{
		ClassName:        request.ClassName,
		ClassRequirement: request.ClassRequirement,
		ClassDateTime:    request.ClassDateTime,
		ClassCapacity:    request.ClassCapacity,
		TrainerID:        request.TrainerID,
	}
	if err := h.ClassService.CreateTrainingClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create class"})
		return
	}
	utils.SendSuccessResponse(c, "add classes Succesfull", class)
}

func (h *ClassHandler) GetClasses(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	date := c.DefaultQuery("date", "")

	classes, err := h.ClassService.GetClasses(id, date)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get classes")
		return
	}
	utils.SendSuccessResponse(c, "Classes found", classes)
}

func (h *ClassHandler) GetClassesHistory(c *gin.Context) {
	idUserStr := c.DefaultQuery("idUser", "")
	idUser, err := strconv.ParseUint(idUserStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "invalid id")
		return
	}
	idUserUint := uint(idUser)

	classes, err := h.ClassService.GetClassesHistory(idUserUint)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user history")
		return
	}
	utils.SendSuccessResponse(c, "Classes found", classes)
}

func (h *ClassHandler) BookClass(c *gin.Context) {
	// Step 1: Parse request body
	var request BookClassRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	// Step 2: Retrieve class details
	classes, err := h.ClassService.GetClasses(fmt.Sprintf("%d", request.IDClass), "")
	if err != nil || len(classes) == 0 {
		utils.SendErrorResponse(c, http.StatusNotFound, "Class not found")
		return
	}
	class := classes[0]

	// Step 3: Check if the class is full
	totalParticipants, err := h.ClassService.CountParticipant(request.IDClass)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to count participants")
		return
	}
	if class.ClassCapacity <= totalParticipants {
		utils.SendErrorResponse(c, http.StatusConflict, "Class is full")
		return
	}

	// Step 4: Check if the user has already booked the class
	alreadyBooked, err := h.ClassService.AlreadyBooked(request.IDUser, class.IDClass)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error checking booking status")
		return
	}
	if alreadyBooked {
		utils.SendErrorResponse(c, http.StatusConflict, "User has already booked this class")
		return
	}

	// Step 5: Verify user's membership/visit for the class date
	if request.Type == "membership" {
		hasMembership, err := h.verifyMembership(request.IDUser, class.ClassDate)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Membership verification failed: %v", err))
			return
		}
		if !hasMembership {
			utils.SendErrorResponse(c, http.StatusForbidden, "User does not have an active or verified membership for this class date")
			return
		}
		trainingClassDetail := models.TrainingClassDetail{
			UserID:          &request.IDUser,
			TrainingClassID: &class.IDClass,
		}
		if err := h.ClassService.BookClass(&trainingClassDetail); err != nil {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to book the class")
			return
		}
		utils.SendSuccessResponse(c, "Successfully booked class", nil)

	} else if request.Type == "visit" {
		visits, err := h.VisitService.GetVisitByUser(request.IDUser)
		if len(visits) == 0 {
			utils.SendErrorResponse(c, http.StatusNotFound, "User has no visit packages")
			return
		}
		if err != nil {
			log.Printf("Error retrieving user visit ")
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get visits")
			return
		}
		var visitAvailable bool = false
		for _, visit := range visits {
			if visit.VisitNumber > visit.VisitUsed &&
				visit.PaymentStatus == "VERIFIED" {
				visitAvailable = true

				err := h.VisitService.UseVisit(request.IDUser, visit.VisitID)
				if err != nil {
					utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to use visit")
					return
				}
				break
			}
		}
		if !visitAvailable {
			utils.SendErrorResponse(c, http.StatusNotFound, "User has no visit packages available or visit package payment not verified")
			return
		}

		trainingClassDetail := models.TrainingClassDetail{
			UserID:          &request.IDUser,
			TrainingClassID: &class.IDClass,
		}
		if err := h.ClassService.BookClass(&trainingClassDetail); err != nil {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to book the class")
			return
		}
		utils.SendSuccessResponse(c, "Successfully booked class", nil)
	} else {
		utils.SendErrorResponse(c, http.StatusBadRequest, "type is empty ")
	}
}

// verifyMembership checks if a user has a valid membership for the given class date
func (h *ClassHandler) verifyMembership(userID uint, classDate time.Time) (bool, error) {
	memberships, err := h.MembershipService.GetMembershipByUser(userID)
	if len(memberships) == 0 {
		return false, fmt.Errorf("no memberships found for the user")
	}
	if err != nil {
		log.Printf("Error retrieving memberships for user %d: %v", userID, err)
		return false, fmt.Errorf("unable to retrieve memberships")
	}

	// Verify active and verified memberships
	for _, membership := range memberships {
		if membership.StartDate.Before(classDate) &&
			membership.EndDate.After(classDate) &&
			membership.PaymentStatus == "VERIFIED" {
			return true, nil
		}
	}

	return false, nil
}
func (h *ClassHandler) GetTrainerSchedule(c *gin.Context) {
	idTrainerStr := c.DefaultQuery("id", "")
	idTrainer, err := strconv.ParseUint(idTrainerStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "invalid id")
		return
	}
	idTrainerUint := uint(idTrainer)
	schedule, err := h.ClassService.GetTrainerSchedule(idTrainerUint)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get trainer schedule")
		return
	}

	utils.SendSuccessResponse(c, "trainer schedule found", schedule)

}
func (h *ClassHandler) CheckEligibility(c *gin.Context) {
	response := EligibilityResponse{
		ValidMember:   true,
		ValidVisit:    true,
		AlreadyBooked: false,
	}

	iduserstr := c.DefaultQuery("iduser", "")
	iduser, err := strconv.ParseUint(iduserstr, 10, 32)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}
	iduserUint := uint(iduser)

	idclass := c.DefaultQuery("idclass", "")
	classes, err := h.ClassService.GetClasses(idclass, "")
	if err != nil || len(classes) == 0 {
		log.Printf("Failed to retrieve classes: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve classes")
		return
	}

	// Step 4: Check if the user has already booked the class
	alreadyBooked, err := h.ClassService.AlreadyBooked(iduserUint, classes[0].IDClass)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error checking booking status")
		return
	}
	if alreadyBooked {
		fmt.Println(alreadyBooked)
		response.AlreadyBooked = true
	}

	hasMember, err := h.verifyMembership(iduserUint, classes[0].ClassDate)
	if err != nil {
		response.ValidMember = false
	}
	if !hasMember {
		response.ValidMember = false
	}

	visits, err := h.VisitService.GetVisitByUser(iduserUint)
	if err != nil {
		response.ValidVisit = false

	}

	response.ValidVisit = false
	for _, visit := range visits {
		if visit.VisitNumber > visit.VisitUsed && visit.PaymentStatus == "VERIFIED" {
			response.ValidVisit = true
			break
		}
	}

	utils.SendSuccessResponse(c, "Eligibility check completed", response)
}
