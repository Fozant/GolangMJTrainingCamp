package routes

import (
	authController "GoMJTrainingCamp/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	authRoutes(r)

}

// authRoutes defines authentication-related routes
func authRoutes(r *gin.Engine) {
	r.GET("/api/login", authController.HandleLogin)
	r.POST("/api/register", authController.HandleRegister)
}
