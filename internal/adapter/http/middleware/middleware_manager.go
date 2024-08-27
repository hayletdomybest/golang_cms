package middleware

import (
	"net/http"
	"strings"
	"weex_admin/internal/context"
	"weex_admin/internal/domain/validators"
	"weex_admin/internal/domain/valueobjects"
	"weex_admin/internal/infrastructure/jwt"

	"github.com/gin-gonic/gin"
)

type MiddlewareManager struct {
	jwtManager    *jwt.JWTManager
	authValidator validators.AuthValidator
}

func NewMiddlewareManager(jwtManager *jwt.JWTManager, authValidator validators.AuthValidator) *MiddlewareManager {
	return &MiddlewareManager{
		jwtManager:    jwtManager,
		authValidator: authValidator,
	}
}

func (manager *MiddlewareManager) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Request does not contain an access token"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		payload, err := manager.jwtManager.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_context", &context.UserContext{
			UserID:   payload.UserID,
			UserName: payload.UserName,
			RoleName: payload.RoleName,
		})

		c.Next()
	}
}

func (manager *MiddlewareManager) CheckPermission(policy valueobjects.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := c.MustGet("user_context").(*context.UserContext)

		pass, err := manager.authValidator.CheckPermission(&validators.CheckPermissionRequest{
			UserID:   &userContext.UserID,
			UserName: &userContext.UserName,
			Policy:   policy,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !pass {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
