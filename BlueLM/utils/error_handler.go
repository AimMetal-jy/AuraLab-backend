package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Error represents a standard error response.
// swagger:response errorResponse
type Error struct {
	// The error message
	// in: body
	Message string `json:"message"`
}

// ErrorHandler is a middleware to handle errors gracefully.
func ErrorHandler(c *gin.Context, err error, statusCode int, message string) {
	logrus.WithError(err).Error(message)
	c.JSON(statusCode, Error{Message: message})
}

// AbortWithInternalServerError responds with a 500 Internal Server Error.
func AbortWithInternalServerError(c *gin.Context, err error) {
	ErrorHandler(c, err, http.StatusInternalServerError, "An unexpected error occurred. Please try again later.")
}

// AbortWithBadRequest responds with a 400 Bad Request error.
func AbortWithBadRequest(c *gin.Context, err error, message string) {
	ErrorHandler(c, err, http.StatusBadRequest, message)
}