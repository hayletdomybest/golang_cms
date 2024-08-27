package user

import (
	"net/http"
	"strconv"
	"weex_admin/internal/application/services"
	"weex_admin/internal/infrastructure/jwt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService *services.AuthService
	jwtManager  *jwt.JWTManager

	resource string
}

// NewUserController creates a new UserController
func NewUserController(
	authService *services.AuthService,
	jwtManager *jwt.JWTManager,
) *UserController {
	return &UserController{
		authService: authService,
		jwtManager:  jwtManager,
	}
}

// CreateUser handles the creation of a new user
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleName string `json:"role_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := ctrl.authService.CreateUser(&services.CreateUserRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		RoleName: &input.RoleName,
		Tx:       nil},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, struct {
			ID       uint64 `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			RoleName string `json:"role_name"`
		}{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleName: user.Role.Name,
		})
	}
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleName string `json:"role_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := ctrl.authService.CreateUser(&services.CreateUserRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		RoleName: &input.RoleName,
		Tx:       nil},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, struct {
			ID       uint64 `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			RoleName string `json:"role_name"`
		}{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			RoleName: user.Role.Name,
		})
	}
}

// GetUserByID handles retrieving a user by ID
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := ctrl.authService.GetUserByID(nil, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// More handlers (UpdateUser, DeleteUser) can be added here
