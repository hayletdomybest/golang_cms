package bootstrap

import (
	"log"
	"weex_admin/internal/application/services"
	"weex_admin/internal/domain/entities"
	"weex_admin/internal/domain/repositories"
	sharederrors "weex_admin/internal/shared/errors"
	"weex_admin/internal/shared/roles"

	"gorm.io/gorm"
)

type DatabaseInitializer struct {
	db              *gorm.DB
	dbTransactioner repositories.DBTransactioner
	authService     *services.AuthService
	authRepository  repositories.AuthRepository
}

func NewDatabaseInitializer(
	db *gorm.DB,
	dbTransactioner repositories.DBTransactioner,
	authService *services.AuthService,
	authRepository repositories.AuthRepository,
) *DatabaseInitializer {
	return &DatabaseInitializer{
		db:              db,
		dbTransactioner: dbTransactioner,
		authService:     authService,
		authRepository:  authRepository,
	}
}

func (initializer *DatabaseInitializer) SeedDatabase() {
	db := initializer.db
	db.AutoMigrate(&entities.Role{}, &entities.User{})
	tx, err := initializer.dbTransactioner.Begin()
	if err != nil {
		log.Fatalf("failed to start transaction: %v", err)
	}
	initializer.ensureRole(tx.(*gorm.DB))
	initializer.ensureAdminUser(tx.(*gorm.DB))
	if err := initializer.dbTransactioner.Commit(tx); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}
}

func (initializer *DatabaseInitializer) ensureRole(tx *gorm.DB) {
	internPanic := func(e error) {
		log.Fatalf("failed to seed roles: %v", e)
	}
	for _, roleName := range []string{roles.Admin, roles.Basic} {
		if err := initializer.authService.CreateRole(&services.CreateRoleRequest{
			RoleName: roleName,
			Tx:       tx,
		}); err != nil {
			internPanic(err)
		}
	}
}

func (initializer *DatabaseInitializer) ensureAdminUser(tx *gorm.DB) {
	internPanic := func(e error) {
		log.Fatalf("failed to add admin user: %v", e)
	}

	if _, err := initializer.authService.CreateUser(&services.CreateUserRequest{
		Name:     roles.Admin,
		Password: roles.Admin,
		RoleName: &roles.Admin,
		Tx:       tx,
	}); err != nil {
		if !sharederrors.ErrAlreadyExist.Is(err) {
			internPanic(err)
		}
	}
}
