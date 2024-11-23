package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
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

func NewClassHandler(classService service.ClassServiceInterface) *ClassHandler {
	return &ClassHandler{ClassService: classService}
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
	err = h.ClassService.BookClass(&trainingClassDetail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to book the class",
		})
		return
	}
	memberships, err := h.MembershipService.GetMembershipByUser(request.IDUser)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch memberships")
		return
	}

	if len(memberships) == 0 {
		utils.SendErrorResponse(c, http.StatusNotFound, "No memberships found for the user")
		return
	}
	//memberships[0].

	utils.SendSuccessResponse(c, "successfully booked class", nil)
}
