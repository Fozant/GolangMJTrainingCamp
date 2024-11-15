package authController

import (
	trainer2 "GoMJTrainingCamp/dbs/models/trainer"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddTrainerRequest struct {
	TrainerName        string `json:"trainerName"`
	TrainerDescription string `json:"trainerDescription"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	PNumber            string `json:"pNumber"`
}

func AddTrainer(c *gin.Context) {
	var addTrainer AddTrainerRequest
	if err := c.ShouldBindJSON(&addTrainer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trainer := trainer2.Trainer{
		TrainerName:        addTrainer.TrainerName,
		TrainerDescription: addTrainer.TrainerDescription,
	}
	idTrainer, err := service.AddTrainer(&trainer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create class"})
		return
	}
	HandleRegisterTrainer(idTrainer, addTrainer, c)

	utils.SendSuccessResponse(c, "add trainer succesfull", trainer)
}

//
//func AddTrainer(c *gin.Context) {
//	var addTrainer AddTrainerRequest
//	if err := c.ShouldBindJSON(&addTrainer); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	trainer := trainer2.Trainer{
//		TrainerName:        addTrainer.TrainerName,
//		TrainerDescription: addTrainer.TrainerDescription,
//	}
//	idTrainer, err := service.AddTrainer(&trainer)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trainer"})
//		return
//	}
//
//	// Teruskan data request yang sudah dibaca
//	HandleRegisterTrainer(idTrainer, addTrainer, c)
//
//	utils.SendSuccessResponse(c, "add trainer successful", trainer)
//}
