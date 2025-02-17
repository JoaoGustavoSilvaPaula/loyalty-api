package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaogustavosp/loyalty-api/internal/services"
)

type ExtractHandler struct {
	extractService *services.ExtractService
}

func NewExtractHandler(extractService *services.ExtractService) *ExtractHandler {
	return &ExtractHandler{extractService}
}

func (h *ExtractHandler) ExtractData(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "URL not provided", "data": nil})
		return
	}

	products, invoice, err := h.extractService.ExtractData(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error extracting data: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data extracted successfully", "data": gin.H{"invoice": invoice, "products": products}})
}
