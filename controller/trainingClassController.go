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

//func GetClasses(c *gin.Context) {
//	id := c.Query("id")
//	date := c.Query("date")
//
//	var classes []models.TrainingClass
//
//	// Check if 'id' parameter is provided
//	if id != "" {
//		err := dbConnection.DB.Where("id = ?", id).Find(&classes).Error
//		if err != nil {
//			utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
//			return
//		}
//	} else if date != "" {
//		// Check if 'date' parameter is provided
//		err := dbConnection.DB.Where("Date = ?", date).Find(&classes).Error
//		if err != nil {
//			utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
//			return
//		}
//	} else {
//		// No parameter, get all classes
//		err := dbConnection.DB.Find(&classes).Error
//		if err != nil {
//			utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
//			return
//		}
//	}
//	utils.SendSuccessResponse(c, "Classes found", classes)
//}
