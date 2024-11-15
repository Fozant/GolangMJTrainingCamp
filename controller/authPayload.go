package authController

import models "GoMJTrainingCamp/dbs/models/users"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	PNumber  string      `json:"p_number"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     models.Role `json:"role"`
}
type LoginResponse struct {
	token string `json:"token"`
}
