package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, authController *AuthController) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		// Add more auth routes here if needed
	}
}
