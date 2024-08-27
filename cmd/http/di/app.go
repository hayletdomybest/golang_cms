package di

import (
	"weex_admin/internal/application/services"
	"weex_admin/internal/bootstrap"
	"weex_admin/internal/context"
	"weex_admin/internal/domain/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	AppContext          *context.AppContext
	Engine              *gin.Engine
	DatabaseInitializer *bootstrap.DatabaseInitializer
	DB                  *gorm.DB

	AuthRepository repositories.AuthRepository
	AuthService    *services.AuthService

	UserRepository repositories.UserRepository
}

func NewApp(
	appContext *context.AppContext,
	engine *gin.Engine,
	databaseInitializer *bootstrap.DatabaseInitializer,
	db *gorm.DB,
	authRepository repositories.AuthRepository,
	authService *services.AuthService,
	userRepository repositories.UserRepository,
) *App {
	return &App{
		AppContext:          appContext,
		Engine:              engine,
		DatabaseInitializer: databaseInitializer,
		DB:                  db,

		AuthRepository: authRepository,
		AuthService:    authService,
		UserRepository: userRepository,
	}
}
