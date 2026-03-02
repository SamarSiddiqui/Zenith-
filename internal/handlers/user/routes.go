package user

import (
	"habit-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes initializes the routes for user-related operations
func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		userGroup.POST("/signup", Signup)
		userGroup.POST("/login", Login)

		// Protected routes
		protected := userGroup.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", GetProfile)
			protected.PATCH("/profile", PatchProfile)
		}
	}
}
