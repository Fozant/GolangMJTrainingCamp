package authController

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

func CreateClass(c *gin.Context) {
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
	if err := service.CreateTrainingClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create class"})
	}
	utils.SendSuccessResponse(c, "add Product Succesfull", class)
}
