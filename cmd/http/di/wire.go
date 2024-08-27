//go:build wireinject
// +build wireinject

package di

import (
	"weex_admin/internal/adapter/http"
	"weex_admin/internal/adapter/http/auth"
	"weex_admin/internal/adapter/http/middleware"
	"weex_admin/internal/adapter/http/user"
	"weex_admin/internal/application/services"
	"weex_admin/internal/bootstrap"
	"weex_admin/internal/context"
	"weex_admin/internal/domain/validators"
	"weex_admin/internal/infrastructure/database"
	"weex_admin/internal/infrastructure/jwt"
	"weex_admin/internal/infrastructure/repositories/gorm"

	"github.com/google/wire"
)

func InitializeApp(
	appContext *context.AppContext,
	databaseConf *database.DatabaseConfig,
	jwtConf *jwt.JWTManagerConf,
) (*App, error) {
	wire.Build(
		// gorm
		database.NewDatabase,
		gorm.NewDBTransaction,
		// repositories
		gorm.NewRoleRepository,
		gorm.NewUserRepository,
		gorm.NewAuthRepository,

		// jwt
		jwt.NewJWTManager,

		// auth
		wire.Bind(new(validators.AuthValidator), new(*services.AuthService)),
		services.NewAuthService,
		auth.NewAuthController,

		// User
		user.NewUserController,

		middleware.NewMiddlewareManager,
		bootstrap.NewDatabaseInitializer,
		http.NewV1Router,
		NewApp,
	)

	return &App{}, nil
}
