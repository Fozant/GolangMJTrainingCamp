package authController

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClassPayload struct {
}

func CreateClass(c *gin.Context) {

	var class models.TrainingClass

	if err := c.ShouldBindJSON(&class); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	dbConnection.DB.Create(&class)
	c.JSON(http.StatusOK, gin.H{"product": class})
}
