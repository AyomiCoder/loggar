package handlers

import (
	"net/http"

	"github.com/AyomiCoder/loggar/pkg/ai"
	"github.com/gin-gonic/gin"
)

type AnalyzeRequest struct {
	Logs string `json:"logs" binding:"required"`
}

// AnalyzeHandler handles log analysis requests
func AnalyzeHandler(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "logs field is required"})
		return
	}

	// Analyze logs using AI
	result, err := ai.AnalyzeLogs(req.Logs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
