package service

import (
	models2 "GoMJTrainingCamp/dbs/dbConnection"
	//"GoMJTrainingCamp/dbs/models/trainer"
	"GoMJTrainingCamp/dbs/models/users"
	"GoMJTrainingCamp/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	PNumber  string      `json:"p_number"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password" `
	Role     models.Role `json:"role" `
	//IDTrainer *trainer.Trainer `json:"id_trainer,omitempty" gorm:"foreignKey:IDTrainer"`
}

func HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate credentials directly with GetUserByID in dbs package.
	user, err := ValidateUserCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := models.CreateJWT([]byte("my-secret-key"), int(user.IDUser))
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

	// Validate the request struct using the validator
	if err := validate.Struct(request); err != nil {
		// If validation fails, return a detailed error message
		errors := err.(validator.ValidationErrors)
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid payload: %v", errors))
		return
	}

	// Check if the email already exists in the database
	_, err := models.GetUserByEmail(request.Email)
	if err == nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("User with email %s already exists", request.Email))
		return
	}

	// Hash the user's password
	hashedPassword, err := models.HashPassword(request.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Create the new user object
	user := models.User{
		PNumber:  request.PNumber,
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		Role:     request.Role,
		//IDTrainer: request.IDTrainer,
	}

	// Create the user in the database
	err = models.CreateUser(user)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Respond with success
	utils.SendSuccessResponse(c, "User registered successfully", nil)
}

func ValidateUserCredentials(email, password string) (*models.User, error) {
	var user models.User
	// Query the user by email
	if err := models2.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// If user is found and password is correct, return the user object
	return &user, nil
}

// GetUserIDFromContext retrieves the users ID from the context if it exists
func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(models.UserKey).(int)
	if !ok {
		log.Println("User ID not found in context")
		return -1
	}

	return userID
}
