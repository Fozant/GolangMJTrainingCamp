package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"time"
)

type CreateClassRequest struct {
	ClassName        string    `json:"className" binding:"required"`
	ClassRequirement string    `json:"classRequirement,omitempty"`
	ClassDateTime    time.Time `json:"ClassDateTime" binding:"required"`
	ClassCapacity    uint      `json:"classCapacity" binding:"required"`
}

type BookClassRequest struct {
	IDClass uint   `json:"idClass" binding:"required"`
	IDUser  uint   `json:"idUser,omitempty"`
	Type    string `json:"type" binding:"required"`
}

type ClassHandler struct {
	ClassService      service.ClassServiceInterface
	MembershipService service.MembershipServiceInterface
}

func NewClassHandler(classService service.ClassServiceInterface, membershipService service.MembershipServiceInterface) *ClassHandler {
	return &ClassHandler{ClassService: classService,
		MembershipService: membershipService}
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
	}
	if err := h.ClassService.CreateTrainingClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create class"})
	}
	utils.SendSuccessResponse(c, "add classes Succesfull", class)
}

func (h *ClassHandler) GetClasses(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	date := c.DefaultQuery("date", "")
	classes, err := h.ClassService.GetClasses(id, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": err.Error(),
		})
		return
	}
	utils.SendSuccessResponse(c, "Classes found", classes)
}

func (h *ClassHandler) BookClass(c *gin.Context) {
	var request BookClassRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	classes, err := h.ClassService.GetClasses(fmt.Sprintf("%d", request.IDClass), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	total, err := h.ClassService.CountParticipant(request.IDClass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to count participants",
		})
		return
	}
	if classes[0].ClassCapacity <= total {
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Class is full",
		})
		return
	}

	exists, err := h.ClassService.AlreadyBooked(request.IDUser, classes[0].IDTrainingClass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error checking booking status",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "User has already booked this class",
		})
		return
	}

	trainingClassDetail := models.TrainingClassDetail{
		UserID:          &request.IDUser,
		TrainingClassID: &classes[0].IDTrainingClass,
	}

	result, err := h.MembershipService.GetMembershipByUser(request.IDUser)
	if err != nil {
		log.Printf("Error getting user membership: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error getting user membership: %v", err))
		return
	}

	if len(result) == 0 {
		utils.SendErrorResponse(c, http.StatusNotFound, "No memberships found for the user")
		return
	}
	classDate := classes[0].ClassDateTime

	hasActiveMembership := false
	for _, entry := range result {
		fmt.Println(entry.PaymentStatus)
		fmt.Println("Aaaa" + entry.PaymentStatus)
		if entry.StartDate.Before(classDate) &&
			entry.EndDate.After(classDate) &&
			entry.PaymentStatus == "VERIFIED" {

			hasActiveMembership = true
			break
		}
	}

	if !hasActiveMembership {
		utils.SendErrorResponse(c, http.StatusForbidden, "User does not have an active membership or membership is not verified for this class date.")
		return
	}
	err = h.ClassService.BookClass(&trainingClassDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to book the class",
		})
		return
	}
	utils.SendSuccessResponse(c, "Successfully booked class", nil)

}
