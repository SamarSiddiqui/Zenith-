package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse standardizes error responses across the application
func ErrorResponse(c *gin.Context, code int, message string, err error) {
	// Default error messages for HTTP status codes
	defaultMessages := map[int]string{
		http.StatusBadRequest:          "Bad Request",
		http.StatusUnauthorized:        "Unauthorized",
		http.StatusForbidden:           "Forbidden",
		http.StatusNotFound:            "Not Found",
		http.StatusConflict:            "Conflict",
		http.StatusUnprocessableEntity: "Unprocessable Entity",
		http.StatusInternalServerError: "Internal Server Error",
		http.StatusServiceUnavailable:  "Service Unavailable",
	}

	// Fallback to default message if custom message is not provided
	if message == "" {
		if defaultMsg, exists := defaultMessages[code]; exists {
			message = defaultMsg
		} else {
			message = "An error occurred"
		}
	}

	// Build the response payload
	response := gin.H{
		"status":  "error",
		"message": message,
		"code":    code,
	}

	if err != nil {
		response["details"] = err.Error() // Include error details in response
	}

	// Send JSON response
	c.JSON(code, response)
}
