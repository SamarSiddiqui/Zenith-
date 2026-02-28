package user

import (
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes initializes the routes for user-related operations
func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		userGroup.POST("/signup", Signup)
		userGroup.POST("/login", Login)
	}
}
