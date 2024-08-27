package repositories

import (
	"weex_admin/internal/domain/valueobjects"
)

type CheckPermissionRequest struct {
	UserID   *uint64
	UserName *string
	Policy   valueobjects.Policy
}

type AddPoliciesForRoleRequest struct {
	RoleName string
	Policies []valueobjects.Policy
}

type GetRolePolicesRequest struct {
	RoleName string
}

type AddRoleForUserRequest struct {
	UserName string
	RoleName string
	Tx       interface{}
}

type CheckPermissionOption = TransactionOption
type AddPoliciesForRoleOption = TransactionOption
type GetRolePolicesOption = TransactionOption
type AddRoleForUserOption = TransactionOption

type AuthRepository interface {
	CheckPermission(body *CheckPermissionRequest, option *CheckPermissionOption) (bool, error)
	AddPoliciesForRole(body *AddPoliciesForRoleRequest, option *AddPoliciesForRoleOption) error
	GetRolePolices(body *GetRolePolicesRequest, option *GetRolePolicesOption) ([]valueobjects.Policy, error)
	AddRoleForUser(body *AddRoleForUserRequest, option *AddRoleForUserOption) error
}
