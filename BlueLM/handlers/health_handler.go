package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// HealthHandler handles health check requests
func HealthHandler(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Message:   "BlueLM service is running",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}