package gorm

import (
	"weex_admin/internal/domain/entities"
	"weex_admin/internal/domain/repositories"
	"weex_admin/internal/shared/errors"
	sharedgorm "weex_admin/internal/shared/gorm"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func (r *RoleRepository) AddRole(body *repositories.AddRoleRequest, option *repositories.TransactionOption) (*entities.Role, error) {
	tx, err := r.getTx(option)
	if err != nil {
		return nil, err
	}

	role := &entities.Role{Name: body.RoleName}
	if err := tx.Create(role).Error; err != nil {
		return nil, sharedgorm.TryWrapErr(err)
	}

	return role, nil
}

func (r *RoleRepository) GetRole(body *repositories.GetRoleRequest, option *repositories.GetRoleOption) (*entities.Role, error) {
	tx, err := r.getTx(option)
	if err != nil {
		return nil, err
	}

	role := &entities.Role{}
	if body.RoleID != nil {
		role.ID = *body.RoleID
	} else if body.RoleName != nil {
		role.Name = *body.RoleName
	} else {
		return nil, errors.ErrInvalidParams.WithMessage("should pass either role id or role name")
	}

	if err := tx.First(role).Error; err != nil {
		if sharedgorm.IsNotFoundErr(err) {
			return nil, nil
		}
		return nil, sharedgorm.TryWrapErr(err)
	}

	return role, nil
}

func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &RoleRepository{db: db}
}

var _ repositories.RoleRepository = (*RoleRepository)(nil)

func (r *RoleRepository) getTx(option *repositories.TransactionOption) (*gorm.DB, error) {
	db := r.db
	if option != nil && option.Tx != nil {
		if tdDB, ok := option.Tx.(*gorm.DB); !ok {
			return nil, errors.ErrPanic.WithMessagef("should pass gorm.DB as transaction")
		} else {
			db = tdDB
		}
	}
	return db.Model(&entities.Role{}), nil
}
