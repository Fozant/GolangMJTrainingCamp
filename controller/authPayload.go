package authController

import models "GoMJTrainingCamp/dbs/models/users"

// LoginRequest is the DTO for login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest is the DTO for user registration request
type RegisterRequest struct {
	PNumber  string      `json:"p_number"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     models.Role `json:"role"` // You can map this to your Role type later in service
}
type LoginResponse struct {
	token string `json:"token"`
}
