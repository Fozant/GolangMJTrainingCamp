package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type BuyVisitRequest struct {
	UserId           uint   `json:"idUser"`
	Price            uint   `json:"price"`
	VisitNumber      uint   `json:"visitNumber"`
	TransactionPrice uint   `json:"transactionPrice"`
	PaymentMethod    string `json:"paymentMethod"`
}

type VisitHandler struct {
	VisitService       service.VisitServiceInterface
	TransactionService service.TransactionServiceInterface
}

func NewVisitHandler(visitService service.VisitServiceInterface, transactionService service.TransactionServiceInterface) *VisitHandler {
	return &VisitHandler{VisitService: visitService,
		TransactionService: transactionService}
}
func (h *VisitHandler) BuyVisit(c *gin.Context) {
	var request BuyVisitRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Request Body:", c.Request.Body)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	visit := models.VisitPackage{
		UserID:      request.UserId,
		Price:       request.Price,
		VisitUsed:   0,
		VisitNumber: request.VisitNumber,
	}

	visitID, err := h.VisitService.BuyVisit(&visit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	transaction := models.Transaction{
		VisitPackageID:   &visitID,
		PaymentType:      "Visit Package",
		PaymentMethod:    request.PaymentMethod,
		TransactionPrice: request.TransactionPrice,
	}
	if err := h.TransactionService.CreateTransaction(&transaction); err != nil {

		log.Printf("Error creating transaction: %v", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create transaction"})
		return
	}
	utils.SendSuccessResponse(c, "Successfully bought visit package", nil)
}
