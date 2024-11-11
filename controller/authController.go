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

	// Generate JWT token
	token, err := service.CreateJWT([]byte("my-secret-key"), int(user.IDUser))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Send JWT response
	utils.SendSuccessResponse(c, "Login successful", gin.H{"token": token})
}

func HandleRegister(c *gin.Context) {
	var request RegisterRequest

	// Parse the incoming JSON request into the RegisterRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the request struct
	if err := validate.Struct(request); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid payload: %v", errors))
		return
	}

	// Check if the email already exists
	_, err := service.GetUserByEmail(request.Email)
	if err == nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("User with email %s already exists", request.Email))
		return
	}

	// Hash password and create the user
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

	// Create user in the database
	if err := service.CreateUser(user); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Respond with success
	utils.SendSuccessResponse(c, "User registered successfully", nil)
}
