package authController

import (
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
	utils.SendSuccessResponse(c, "Login successful", gin.H{"token": token})
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
	user := models.User{
		PNumber:          request.PNumber,
		Name:             request.Name,
		Email:            request.Email,
		Password:         hashedPassword,
		Role:             models.RoleUser,
		RegistrationDate: time.Now(),
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

func HandleRegisterTrainer(id uint, addTrainer AddTrainerRequest, c *gin.Context) {

	validate := validator.New()

	if err := validate.Struct(addTrainer); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid payload: %v", errors))
		return
	}
	_, err := service.GetUserByEmail(addTrainer.Email)
	if err == nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("User with email %s already exists", addTrainer.Email))
		return
	}
	hashedPassword, err := service.HashPassword(addTrainer.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user := models.User{
		PNumber:          addTrainer.PNumber,
		Name:             addTrainer.Name,
		Email:            addTrainer.Email,
		Password:         hashedPassword,
		Role:             models.RoleTrainer,
		RegistrationDate: time.Now(),
		IDTrainer:        id,
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
