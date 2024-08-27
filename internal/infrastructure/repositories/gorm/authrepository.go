package gorm

import (
	"weex_admin/internal/domain/repositories"
	"weex_admin/internal/domain/valueobjects"
	"weex_admin/internal/infrastructure/casbin"
	"weex_admin/internal/shared/errors"

	casbinv2 "github.com/casbin/casbin/v2"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func (a *AuthRepository) CheckPermission(body *repositories.CheckPermissionRequest, option *repositories.CheckPermissionOption) (bool, error) {
	tx, err := a.getTx(option)
	if err != nil {
		return false, err
	}

	for _, action := range body.Policy.Action.GetActions() {
		if ok, err := tx.Enforce(body.UserName, body.Policy.Resource, string(action)); err != nil {
			return false, errors.ErrInternal.
				Wrap(err).
				WithMessagef("failed to Enforce for casbin")
		} else {
			if !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

func (a *AuthRepository) AddPoliciesForRole(body *repositories.AddPoliciesForRoleRequest, option *repositories.AddPoliciesForRoleOption) error {
	if len(body.Policies) == 0 {
		return errors.ErrInvalidParams.WithMessage("policies should not be empty")
	}
	if len(body.RoleName) == 0 {
		return errors.ErrInvalidParams.WithMessage("role name should not be empty")
	}

	if err := a.addPoliciesForRole(body.RoleName, body.Policies, option); err != nil {
		return err
	}

	return nil
}

func (a *AuthRepository) GetRolePolices(body *repositories.GetRolePolicesRequest, option *repositories.GetRolePolicesOption) ([]valueobjects.Policy, error) {
	enforcer, err := a.getTx(option)
	if err != nil {
		return nil, err
	}
	if len(body.RoleName) == 0 {
		return nil, errors.ErrInvalidParams.WithMessage("role name should not be empty")
	}

	policies, err := enforcer.GetFilteredPolicy(0, body.RoleName)
	if err != nil {
		return nil, errors.ErrInternal.
			Wrap(err).
			WithMessagef("failed to GetFilteredPolicy for casbin")
	}

	var result []valueobjects.Policy
	for _, policy := range policies {
		actions, err := valueobjects.Wrap(policy[2])
		if err != nil {
			return nil, err
		}
		result = append(result, valueobjects.Policy{
			Resource: policy[1],
			Action:   actions[0],
		})
	}
	return result, nil
}

func (a *AuthRepository) AddRoleForUser(body *repositories.AddRoleForUserRequest, option *repositories.AddRoleForUserOption) error {
	if len(body.UserName) == 0 || len(body.RoleName) == 0 {
		return errors.ErrInvalidParams.WithMessage("user name and role name should not be empty")
	}

	tx, err := a.getTx(option)
	if err != nil {
		return err
	}

	_, err = tx.AddRoleForUser(body.UserName, body.RoleName)
	if err != nil {
		return errors.ErrInternal.
			Wrap(err).
			WithMessagef("failed to AddRoleForUser for casbin")
	}
	return nil
}

var _ repositories.AuthRepository = (*AuthRepository)(nil)

func NewAuthRepository(
	db *gorm.DB) repositories.AuthRepository {
	return &AuthRepository{db: db}
}

func (u *AuthRepository) getTx(option *repositories.TransactionOption) (*casbinv2.Enforcer, error) {
	db := u.db
	if option != nil && option.Tx != nil {
		if tdDB, ok := option.Tx.(*gorm.DB); !ok {
			return nil, errors.ErrPanic.WithMessagef("should pass gorm.DB as transaction")
		} else {
			db = tdDB
		}
	}
	return casbin.NewCasbin(db)
}

func (a *AuthRepository) addPoliciesForRole(roleName string, policies []valueobjects.Policy, option *repositories.AddPoliciesForRoleOption) error {
	if len(roleName) == 0 {
		return errors.ErrPanic.WithMessage("role name should not be empty")
	}

	enforcer, err := a.getTx(option)
	if err != nil {
		return err
	}

	for _, p := range policies {
		for _, action := range p.Action.GetActions() {
			if _, err := enforcer.AddPolicy(roleName, p.Resource, action); err != nil {
				return errors.ErrInternal.
					Wrap(err).
					WithMessagef("failed to AddPolicy for casbin")
			}
		}
	}
	return nil
}
