package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CreateClassRequest struct {
	ClassName        string    `json:"className" binding:"required"`
	ClassRequirement string    `json:"classRequirement,omitempty"`
	ClassDateTime    time.Time `json:"ClassDateTime" binding:"required"`
	ClassCapacity    int64     `json:"classCapacity" binding:"required"`
}

type BookClassRequest struct {
	IDClass uint   `json:"idClass" binding:"required"`
	IDUser  uint   `json:"idUser,omitempty"`
	Type    string `json:"type" binding:"required"`
}

type ClassHandler struct {
	ClassService service.ClassServiceInterface
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err,
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

}
