package v1

import (
	"net/http"
	"net/url"

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
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "URL not provided", "data": nil})
		return
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid URL format", "data": nil})
		return
	}

	products, invoice, err := h.extractService.ExtractData(parsedURL.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error extracting data: " + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data extracted successfully", "data": gin.H{"invoice": invoice, "products": products}})
}
