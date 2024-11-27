package controller

import (
	"GoMJTrainingCamp/dbs/models"
	"GoMJTrainingCamp/service"
	"GoMJTrainingCamp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PackageHandler struct {
	PackageService service.PackageServiceInterface
}

func NewPackageHandler(packageService service.PackageServiceInterface) *PackageHandler {
	return &PackageHandler{PackageService: packageService}
}

type PackageListRequest struct {
	PackageName string `json:"package_name" binding:"required"`
	Price       uint   `json:"price"`
	Duration    *uint  `json:"duration"`
	Status      string `json:"status" binding:"required"`
	Type        string `json:"type"`
	VisitNumber *uint  `json:"visit_number" `
}

func (h *PackageHandler) AddPackage(c *gin.Context) {
	var req PackageListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packages := models.PackageList{
		PackageName: req.PackageName,
		Price:       req.Price,
		Duration:    req.Duration,
		Status:      req.Status,
		Type:        req.Type,
		VisitNumber: req.VisitNumber,
	}
	err := h.PackageService.CreatePackage(&packages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SendSuccessResponse(c, "Sucessfuly created package", nil)
}

func (h *PackageHandler) GetPackage(c *gin.Context) {
	id := c.DefaultQuery("id", "")

	packages, err := h.PackageService.GetPackage(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Failed to retrieve packages")
		return
	}
	utils.SendSuccessResponse(c, "Successfully retrieved packages", packages)
}
