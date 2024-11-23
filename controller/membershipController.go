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

type BuyMembershipRequest struct {
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	Price            uint   `json:"price"`
	Duration         uint   `json:"duration"`
	UserId           uint   `json:"idUser"`
	TransactionPrice uint   `json:"transactionPrice"`
	PaymentMethod    string `json:"paymentMethod"`
}

type MembershipHandler struct {
	MembershipService  service.MembershipServiceInterface
	TransactionService service.TransactionServiceInterface
}

func NewMembershipHandler(membershipService service.MembershipServiceInterface, transactionService service.TransactionServiceInterface) *MembershipHandler {
	return &MembershipHandler{MembershipService: membershipService,
		TransactionService: transactionService}
}
func (h *MembershipHandler) BuyMembership(c *gin.Context) {
	var request BuyMembershipRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Request Body:", c.Request.Body)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	startDate, err := utils.ParseFormattedDate(request.StartDate)
	if err != nil {
		fmt.Println("start date:", request.StartDate)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid start date format. Expected yyyy-mm-dd"})
		return
	}
	endDate, err := utils.ParseFormattedDate(request.EndDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid end date format. Expected yyyy-mm-dd"})
		return
	}

	membership := models.Membership{
		UserID:    request.UserId,
		StartDate: startDate,
		EndDate:   endDate,
		Price:     request.Price,
		Duration:  request.Duration,
	}

	membershipID, err := h.MembershipService.BuyMembership(&membership)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	transaction := models.Transaction{
		MembershipID:     &membershipID,
		PaymentType:      "Membership",
		PaymentMethod:    request.PaymentMethod,
		TransactionPrice: request.TransactionPrice,
		PaymentStatus:    "waiting for approval",
	}
	transactionID, err := h.TransactionService.CreateTransaction(&transaction)
	if err != nil {
		log.Printf("Error creating transaction: with id %d,%v", transactionID, err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	//err = h.MembershipService.UpdateTransactionID(membershipID, transactionID)
	//if err != nil {
	//	log.Printf("Error updating transaction ID in membership: %v", err)
	//	utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update transaction ID")
	//	return
	//}

	utils.SendSuccessResponse(c, "Successfully bought membership", nil)
}
