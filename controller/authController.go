package controller

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	models "GoMJTrainingCamp/dbs/models/users"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"time"
)

func HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Validate credentials directly with GetUserByID in dbs package.
	user, err := service.ValidateUserCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := service.CreateJWT([]byte("my-secret-key"), int(user.IDUser))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.SendSuccessResponse(c, "Login successful", gin.H{
		"user_id":   user.IDUser,
		"user_role": user.Role,
		"user_name": user.Name,
		"token":     token,
	})
}

func HandleRegister(c *gin.Context) {
	var request RegisterRequest

	// Parse the incoming JSON request into the RegisterRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid payload: %v", errors))
		return
	}
	_, err := service.GetUserByEmail(request.Email)
	if err == nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("User with email %s already exists", request.Email))
		return
	}
	hashedPassword, err := service.HashPassword(request.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user := &models.User{
		PNumber:          request.PNumber,
		Name:             request.Name,
		Email:            request.Email,
		Password:         hashedPassword,
		Role:             models.RoleUser,
		RegistrationDate: time.Now(),
		IDTrainer:        nil,
	}
	if err := service.CreateUser(user); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error creating user")
		return
	}
	// Generate JWT token
	token, err := service.CreateJWT([]byte("my-secret-key"), int(user.IDUser))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.SendSuccessResponse(c, "User registered successfully", token)
}

func HandleRegisterTrainer(addTrainer AddTrainerRequest, c *gin.Context) (uint, error) {

	validate := validator.New()

	if err := validate.Struct(addTrainer); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid payload: %v", errors))
		return 0, err
	}
	_, err := service.GetUserByEmail(addTrainer.Email)
	if err == nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("User with email %s already exists", addTrainer.Email))
		return 0, err
	}
	hashedPassword, err := service.HashPassword(addTrainer.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error hashing password")
		return 0, err
	}
	user := &models.User{
		PNumber:          addTrainer.PNumber,
		Name:             addTrainer.Name,
		Email:            addTrainer.Email,
		Password:         hashedPassword,
		Role:             models.RoleTrainer,
		RegistrationDate: time.Now(),
		IDTrainer:        nil,
	}
	if err := service.CreateUser(user); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error creating trainer")
		return 0, err
	}

	token, err := service.CreateJWT([]byte("my-secret-key"), int(user.IDUser))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return 0, err
	}
	utils.SendSuccessResponse(c, "Trainer registered successfully",
		gin.H{
			"token": token,
		})

	return user.IDUser, nil
}

func UpdateUserTrainer(idUser uint, idTrainer uint, c *gin.Context) error {

	result := dbConnection.DB.Model(&models.User{}).
		Where("id_user = ?", idUser).
		Update("id_trainer", idTrainer)

	if result.Error != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error updating trainer")
		return result.Error
	}
	return nil
}
