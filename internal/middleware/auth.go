package middleware

import (
	"net/http"
	"strings"

	"habit-tracker/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization is required", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization  token", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err)
			c.Abort()
			return
		}

		// Extract user_id and set it in the context
		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
