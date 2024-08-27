package http

import (
	"weex_admin/internal/adapter/http/auth"
	"weex_admin/internal/adapter/http/middleware"
	"weex_admin/internal/adapter/http/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewV1Router initializes the Gin engine and registers all module routes
func NewV1Router(
	middlewareManager *middleware.MiddlewareManager,
	authController *auth.AuthController,
	userController *user.UserController,
) *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	group := r.Group("/api/v1")

	// Register Auth module routes
	auth.RegisterAuthRoutes(group, authController)

	// Register User module routes
	user.RegisterUserRoutes(group, middlewareManager, userController)

	// Future modules can be registered here

	return r
}
