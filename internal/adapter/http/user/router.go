package user

import (
	"weex_admin/internal/adapter/http/middleware"
	"weex_admin/internal/domain/valueobjects"
	"weex_admin/internal/shared/constants"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, middlewareManager *middleware.MiddlewareManager, userController *UserController) {
	userGroup := r.Group("/user")
	{
		userGroup.Use(middlewareManager.JWTMiddleware())

		userGroup.POST("/", middlewareManager.CheckPermission(valueobjects.Policy{
			Resource: constants.USER_CREATE_PATH,
			Action:   valueobjects.Write,
		}), userController.CreateUser)

		userGroup.GET("/:id", middlewareManager.CheckPermission(valueobjects.Policy{
			Resource: constants.USER_GET_PATH,
			Action:   valueobjects.Read,
		}), userController.GetUserByID)

		// Add more auth routes here if needed
	}
}
