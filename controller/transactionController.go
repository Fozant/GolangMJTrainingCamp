package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VerifyTransactionRequest struct {
	TransactionID     uint   `json:"transactionId" binding:"required"`
	TransactionStatus string `json:"transactionStatus" binding:"required"`
	Notes             string `json:"notes"`
}

type TransactionHandler struct {
	TransactionService service.TransactionServiceInterface
}

func NewTransactionHandler(TransactionService service.TransactionServiceInterface) *TransactionHandler {
	return &TransactionHandler{TransactionService: TransactionService}
}

func (h *TransactionHandler) VerifyTransaction(c *gin.Context) {
	var request VerifyTransactionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	transaction := models.Transaction{
		IDTransaction:     request.TransactionID,
		PaymentStatus:     request.TransactionStatus,
		PaymentStatusNote: request.Notes,
	}

	if err := h.TransactionService.UpdateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}
	utils.SendSuccessResponse(c, "Update transaction sucessfull", transaction)
}
