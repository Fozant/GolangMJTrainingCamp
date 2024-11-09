package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Standard API response structure
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Helper function to send a generic response
func SendResponse(c *gin.Context, status int, message string, data interface{}) {
	response := APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, response)
}

// Utility function to send a success response
func SendSuccessResponse(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusOK, message, data)
}

// Utility function to send an error response
func SendErrorResponse(c *gin.Context, status int, message string) {
	SendResponse(c, status, message, nil)
}
