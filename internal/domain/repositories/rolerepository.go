package repositories

import "weex_admin/internal/domain/entities"

type GetRoleRequest struct {
	RoleID   *uint
	RoleName *string
}
type AddRoleRequest struct {
	RoleName string
}

type GetRoleOption = TransactionOption
type AddRoleOption = TransactionOption

type RoleRepository interface {
	GetRole(body *GetRoleRequest, option *GetRoleOption) (*entities.Role, error)
	AddRole(body *AddRoleRequest, option *AddRoleOption) (*entities.Role, error)
}
