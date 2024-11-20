package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BuyMembershipRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Price     uint   `json:"price"`
	Duration  uint   `json:"duration"`
	UserId    uint   `json:"idUser"`
}

type MembershipHandler struct {
	MembershipService service.MembershipServiceInterface
}

func NewMembershipHandler(membershipService service.MembershipServiceInterface) *MembershipHandler {
	return &MembershipHandler{MembershipService: membershipService}
}
func (h *MembershipHandler) BuyMembership(c *gin.Context) {
	var request BuyMembershipRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Request Body:", c.Request.Body) // Debugging step
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	startDate, err := utils.ParseFormattedDate(request.StartDate)
	if err != nil {
		fmt.Println("start date:", request.StartDate) // Debugging step
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

	if err := h.MembershipService.BuyMembership(&membership); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	utils.SendSuccessResponse(c, "Successfully bought membership", nil)
}
